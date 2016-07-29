#! /usr/bin/env ruby
require 'erb'
require_relative 'utils'

class VATs
  include Utils

  VATS_MANIFEST_ERB_PATH = 'cf-persist-service-broker/ci/templates/vats-manifest.yml.erb'
  VATS_MANIFEST_YML_PATH = 'cf-persist-service-broker/ci/templates/vats-manifest.yml'
  RUN_VATS_SH_PATH = 'cf-persist-service-broker/ci/templates/run_vats.sh'
  REXRAY_CONFIG_ERB_PATH = 'cf-persist-service-broker/ci/templates/rexray_config.json.erb'
  REXRAY_CONFIG_JSON_PATH = 'cf-persist-service-broker/ci/templates/rexray_config.json'

  ENV_PARAMS = %W{
    BOSH_DIRECTOR_PUBLIC_IP
    BOSH_PASSWORD
    BOSH_USER
    REXRAY_RELEASE_NAME
    SCALEIO_ENDPOINT
    SCALEIO_INSECURE
    SCALEIO_MDM_IPS
    SCALEIO_PASSWORD
    SCALEIO_SDC_RELEASE_NAME
    SCALEIO_STORAGE_POOL_NAME
    SCALEIO_USERNAME
    SCALEIO_VERSION
    STORAGE_SERVICE_TYPE
    VATS_DEPLOYMENT_IP
    VATS_DEPLOYMENT_NAME
    VATS_DEPLOYMENT_PASSWORD
  }

  def initialize
    @failed = false
    check_env_params(ENV_PARAMS)
  end

  def perform
    exec_cmd('apt-get -y update && apt-get install -y sshpass')
    exec_cmd('gem install bosh_cli --no-ri --no-rdoc')
    exec_cmd("bosh target #{@bosh_director_public_ip}")
    exec_cmd("bosh login #{@bosh_user} #{@bosh_password}")

    output = exec_cmd('bosh status --uuid')
    bosh_director_uuid = output[1].strip
    generate_bosh_manifest(bosh_director_uuid)

    upload_bosh_releases
    exec_cmd("bosh deployment #{VATS_MANIFEST_YML_PATH}")
    exec_cmd('bosh -n deploy')

    generate_raxray_config
    scp(REXRAY_CONFIG_JSON_PATH)
    scp(RUN_VATS_SH_PATH)
    ssh_run_vats
  ensure
    # exec_cmd("bosh -n delete deployment #{@vats_deployment_name}")
    # exec_cmd("bosh -n delete release #{@rexray_release_name}")
    # exec_cmd("bosh -n delete release #{@scaleio_sdc_release_name}")
  end

  def scp(filepath)
    exec_cmd("sshpass -p #{@vats_deployment_password} \
              scp -o StrictHostKeyChecking=no #{filepath} \
              vcap@#{@vats_deployment_ip}:/home/vcap/")
  end

  def ssh_run_vats
    exec_cmd("sshpass -p #{@vats_deployment_password} \
              ssh -o StrictHostKeyChecking=no vcap@#{@vats_deployment_ip} \
              \"echo #{@vats_deployment_password} | sudo -S bash -c '/home/vcap/run_vats.sh'\"")
  end

  def upload_bosh_releases
    exec_cmd("pushd rexray-bosh-release && \
              bosh -n create release --force --name #{@rexray_release_name} && \
              bosh -n upload release && \
              popd")

    if @storage_service_type == 'scaleio'
      exec_cmd("pushd scaleio-sdc-bosh-release && \
                bosh -n create release --force --name #{@scaleio_sdc_release_name} && \
                bosh -n upload release && \
                popd")
    end
  end

  def generate_bosh_manifest(bosh_director_uuid)
    results = ERB.new(File.read(VATS_MANIFEST_ERB_PATH)).result(binding())
    File.open(VATS_MANIFEST_YML_PATH, 'w+') do |f|
      f.write(results)
    end
  end

  def generate_raxray_config
    results = ERB.new(File.read(REXRAY_CONFIG_ERB_PATH)).result(binding())
    File.open(REXRAY_CONFIG_JSON_PATH, 'w+') do |f|
      f.write(results)
    end
  end

  if __FILE__ == $0
    VATs.new.perform
  end
end
