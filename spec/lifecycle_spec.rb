require File.expand_path '../spec_helper.rb', __FILE__

describe 'lifecycle', type: :lifecycle do
  before(:all) do
    @persist_service_broker_name = ENV['CF_SCALEIO_SB_SERVICE']
    @persist_service_broker_app = 'persist-service-broker-app-ci'
    @persist_service_name = 'persist-service-ci'
    @persist_acceptance_app = 'persist-acceptance-app'
    @uploaded_data = "TADATADA"
    @download_url = "https://#{@persist_acceptance_app}.#{ENV['CF_ENDPOINT']}/data"

    exec_command("cf api https://api.#{ENV['CF_ENDPOINT']} --skip-ssl-validation")
    exec_command("cf auth #{ENV['CF_USERNAME']} \"#{ENV['CF_PASSWORD']}\"")
    exec_command("cf target -o #{ENV['CF_ORG']} -s #{ENV['CF_SPACE']}")
    exec_command("cf push #{@persist_service_broker_app} --no-start")
    set_env
    exec_command("cf start #{@persist_service_broker_app}")
  end

  after(:all) do
    unbind_service
    restage_app(@persist_acceptance_app)
    delete_service_instance
    delete_service_broker
    delete_app(@persist_acceptance_app)
    delete_app(@persist_service_broker_app)
  end

  it 'should register service broker to cf ' do
    get_service_catalog
    register_the_service_broker
    create_service_instance

    test_that_data_persist
  end
end

def exec_command(command)
  puts "Running #{command}"
  output = `#{command}`
  puts output
  output
end

def delete_app(app_name)
  exec_and_check("cf delete #{app_name} -f")
end

def test_that_data_persist
  Dir.chdir('../scaleio-acceptance-app')
  exec_and_check("cf push #{@persist_acceptance_app} --no-start")
  exec_and_check("cf env #{@persist_acceptance_app}")
  bind_service
  exec_and_check("cf start #{@persist_acceptance_app}")

  exec_command("curl --insecure -X POST https://#{@persist_acceptance_app}.#{ENV['CF_ENDPOINT']}/data -d \"#{@uploaded_data}\" -H \"Content-Type: text/plain\"")
  check_persistent(@uploaded_data, @download_url)

  restage_app(@persist_acceptance_app)
  check_persistent(@uploaded_data, @download_url)

  exec_and_check("cf stop #{@persist_acceptance_app}")
  exec_and_check("cf start #{@persist_acceptance_app}")
  check_persistent(@uploaded_data, @download_url)

  unbind_service
  exec_and_check("cf delete #{@persist_acceptance_app} -f")
  exec_and_check("cf push #{@persist_acceptance_app} --no-start")
  bind_service
  exec_and_check("cf start #{@persist_acceptance_app}")
  check_persistent(@uploaded_data, @download_url)
end

def set_env
  exec_command("cf set-env #{@persist_service_broker_app} BROKER_PASSWORD #{ENV['BROKER_PASSWORD']}")
  exec_command("cf set-env #{@persist_service_broker_app} BROKER_USERNAME #{ENV['BROKER_USERNAME']}")
  exec_command("cf set-env #{@persist_service_broker_app} LIBSTORAGE_URI #{ENV['LIBSTORAGE_URI']}")
  exec_command("cf set-env #{@persist_service_broker_app} LIBSTORAGE_STORAGE_DRIVER #{ENV['LIBSTORAGE_STORAGE_DRIVER']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_ENDPOINT #{ENV['SCALEIO_ENDPOINT']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_INSECURE #{ENV['SCALEIO_INSECURE']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_PASSWORD #{ENV['SCALEIO_PASSWORD']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_PROTECTION_DOMAIN_ID #{ENV['SCALEIO_PROTECTION_DOMAIN_ID']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_PROTECTION_DOMAIN_NAME #{ENV['SCALEIO_PROTECTION_DOMAIN_NAME']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_STORAGE_POOL_NAME #{ENV['SCALEIO_STORAGE_POOL_NAME']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_SYSTEM_ID #{ENV['SCALEIO_SYSTEM_ID']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_SYSTEM_NAME #{ENV['SCALEIO_SYSTEM_NAME']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_THIN_OR_THICK #{ENV['SCALEIO_THIN_OR_THICK']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_USE_CERTS #{ENV['SCALEIO_USE_CERTS']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_USERNAME #{ENV['SCALEIO_USERNAME']}")
  exec_command("cf set-env #{@persist_service_broker_app} SCALEIO_VERSION #{ENV['SCALEIO_VERSION']}")
  exec_command("cf env #{@persist_service_broker_app}")
end

def restage_app(app)
  output=exec_command("cf restage #{app}")
  expect(output).to include('OK')
end

def check_persistent(uploaded_data, download_url)
  uri = URI(download_url)
  req = Net::HTTP::Get.new(uri.path)
  res = Net::HTTP.start(
    uri.host, uri.port,
    :use_ssl => uri.scheme == 'https',
    :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |https|

    https.request(req)
  end

  expect(res.body).to eq(uploaded_data)
end

def get_service_catalog
  uri = URI("https://#{@persist_service_broker_app}.#{endpoint}/v2/catalog")
  req = Net::HTTP::Get.new(uri.path)
  req.basic_auth ENV['BROKER_USERNAME'], ENV['BROKER_PASSWORD']
  res = Net::HTTP.start(
    uri.host, uri.port,
    :use_ssl => uri.scheme == 'https',
    :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |https|

    https.request(req)
  end

  expected_catalog = JSON.parse(File.read(File.join(project_path, 'templates/catalog.json')))
  expect(res.code).to eq('200')
  expect(JSON.parse(res.body)).to eq(expected_catalog)
end

def register_the_service_broker
  output = exec_command("cf create-service-broker #{@persist_service_broker_name} #{ENV['BROKER_USERNAME']} #{ENV['BROKER_PASSWORD']} https://#{@persist_service_broker_app}.#{endpoint}")
  expect(output).to include('OK')

  output = exec_command("cf enable-service-access #{@persist_service_broker_name}")
  expect(output).to include('OK')

  output = exec_command('cf marketplace')
  expect(output).to include(@persist_service_broker_name)
  expect(output).to include('ci')

  output = exec_command("cf marketplace -s #{@persist_service_broker_name}")
  expect(output).to include('ci')
  expect(output).to include('free')
end

def create_service_instance
  output = exec_command("cf create-service #{@persist_service_broker_name} ci #{@persist_service_name} -c \'{\"storage_pool_name\": \"default\"}\'")
  expect(output).to include('OK')
end

def bind_service
  output = exec_command("cf bind-service #{@persist_acceptance_app} #{@persist_service_name}")
  expect(output).to include('OK')

  output = exec_command("cf env #{@persist_acceptance_app}")
  expect(output).to include("#{@persist_service_broker_name}")
end

def unbind_service
  output = exec_command("cf unbind-service #{@persist_acceptance_app} #{@persist_service_name}")
  expect(output).to include('OK')

  output = exec_command("cf env #{@persist_acceptance_app}")
  expect(output).to_not include("#{@persist_service_broker_name}")
end

def delete_service_instance
  output = exec_command("cf delete-service #{@persist_service_name} -f")
  expect(output).to include('OK')
end


def delete_service_broker
  output = exec_command("cf delete-service-broker #{@persist_service_broker_name} -f")
  expect(output).to include('OK')
end

def exec_and_check(command)
  output = exec_command(command)
  expect(output).to include('OK')
  output
end
