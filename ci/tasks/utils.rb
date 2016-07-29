require 'open3'
require 'shellwords'

module Utils
  def check_env_params(expected_params)
    expected_params.each do |param|
      puts param
      raise ArgumentError, "'#{param}' not set in ENV" if ENV[param].empty? || ENV[param].eql?("replace-me")
      instance_variable_set("@#{param.downcase}", ENV[param])
    end
  end

  def exec_cmd(command, supress=false)
    cmd = "bash -c #{command.shellescape}"
    puts "---------------"
    puts "Running cmd: #{cmd}:"

    output = ''
    error = ''
    expected_error = nil

    return_stderr = false
    Open3.popen3(cmd) do |stdin, stdout, stderr, wait_thr|
      # Create a thread to read from each stream
      err_thr = Thread.new do
        Thread.current.abort_on_exception = true
        until (line = stderr.gets).nil? do
          error += line
          puts line
        end
      end
      out_thr = Thread.new do
        Thread.current.abort_on_exception = true
        until (line = stdout.gets).nil? do
          output += line
          puts line
        end
      end

      err_thr.join
      out_thr.join
      # Don't exit until the external process is done too
      wait_thr.join

      wait_thr_status = wait_thr.value
      exit_status = wait_thr_status.exitstatus
      termsig = wait_thr_status.termsig
      unless wait_thr_status.success?
        display_output = true
        msg = "Command \"#{cmd}\" failed."
        if wait_thr_status.signaled?
        else
          if expected_error != exit_status
          else
            display_output = false
          end
        end
        if display_output then
          warn "Output was: #{output}"
          warn "Stderr was: #{error}"
        end
        # Note: WOW KEE THIS IS BAD HACK
        if exit_status != 0 && !supress
          @failed = true
          raise
        end

        return [exit_status, output] if not return_stderr
        return [exit_status, output, error]
      end
    end
    return [0, output] if not return_stderr
    return [0, output, error]
  rescue => e
    warn "Command \"#{cmd}\" failed with exception: #{e.message}"
    @failed = true
    raise
    #  return [1, output] if not return_stderr
    #  return [1, output, error]
  end
end
