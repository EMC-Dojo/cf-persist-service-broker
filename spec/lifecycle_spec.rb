require File.expand_path '../spec_helper.rb', __FILE__

describe 'lifecycle', type: :lifecycle do
  before(:all) do
    exec_command("cf api https://api.#{ENV['CF_ENDPOINT']} --skip-ssl-validation")
    exec_command("cf auth #{ENV['CF_USERNAME']} \"#{ENV['CF_PASSWORD']}\"")
    exec_command("cf target -o #{ENV['CF_ORG']} -s #{ENV['CF_SPACE']}")
    exec_command("cf push #{app_name} --no-start -b #{build_pack}")
    exec_command("cf set-env #{app_name} BROKER_USERNAME #{ENV['BROKER_USERNAME']}")
    exec_command("cf set-env #{app_name} BROKER_PASSWORD #{ENV['BROKER_PASSWORD']}")
    exec_command("cf set-env #{app_name} LIBSTORAGE_HOST #{ENV['LIBSTORAGE_HOST']}")
    exec_command("cf set-env #{app_name} LIBSTORAGE_STORAGE_DRIVER #{ENV['LIBSTORAGE_STORAGE_DRIVER']}")
    exec_command("cf set-env #{app_name} SCALEIO_ENDPOINT #{ENV['SCALEIO_ENDPOINT']}")
    exec_command("cf set-env #{app_name} SCALEIO_INSECURE #{ENV['SCALEIO_INSECURE']}")
    exec_command("cf set-env #{app_name} SCALEIO_USE_CERTS #{ENV['SCALEIO_USE_CERTS']}")
    exec_command("cf set-env #{app_name} SCALEIO_USERNAME #{ENV['SCALEIO_USERNAME']}")
    exec_command("cf set-env #{app_name} SCALEIO_PASSWORD #{ENV['SCALEIO_PASSWORD']}")
    exec_command("cf set-env #{app_name} SCALEIO_SYSTEM_ID #{ENV['SCALEIO_SYSTEM_ID']}")
    exec_command("cf set-env #{app_name} SCALEIO_SYSTEM_NAME #{ENV['SCALEIO_SYSTEM_NAME']}")
    exec_command("cf set-env #{app_name} SCALEIO_PROTECTION_DOMAIN_ID #{ENV['SCALEIO_PROTECTION_DOMAIN_ID']}")
    exec_command("cf set-env #{app_name} SCALEIO_PROTECTION_DOMAIN_NAME #{ENV['SCALEIO_PROTECTION_DOMAIN_NAME']}")
    exec_command("cf set-env #{app_name} SCALEIO_STORAGE_POOL_NAME #{ENV['SCALEIO_STORAGE_POOL_NAME']}")
    exec_command("cf set-env #{app_name} SCALEIO_THIN_OR_THICK #{ENV['SCALEIO_THIN_OR_THICK']}")
    exec_command("cf set-env #{app_name} SCALEIO_VERSION #{ENV['SCALEIO_VERSION']}")
    exec_command("cf env #{app_name}")
    exec_command("cf start #{app_name}")
  end

  after(:all) do
    exec_command("cf delete-service-broker scaleio -f")
    exec_command("cf delete #{app_name} -f")
  end

  it 'should push app to cf ' do
    get_service_catalog
    register_the_service_broker
    # create_service_instance
    # bind_service
    # unbind_service
    # delete_service_instance
    # delete_not_created_service
  end
end

def exec_command(command)
  puts "Running #{command}"
  output = `#{command}`
  puts output
  output
end

def get_service_catalog
  uri = URI("https://#{app_name}.#{endpoint}/v2/catalog")
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
  output = exec_command("cf create-service-broker scaleio #{ENV['BROKER_USERNAME']} #{ENV['BROKER_PASSWORD']} https://#{app_name}.#{endpoint}")
  expect(output).to include('OK')

  output = exec_command('cf enable-service-access scaleio')
  expect(output).to include('OK')

  output = exec_command('cf marketplace')
  expect(output).to include('ScaleIO')
  expect(output).to include('small')

  output = exec_command('cf marketplace -s scaleio')
  expect(output).to include('small')
  expect(output).to include('free')
end

def create_service_instance
  output = exec_command('cf create-service scaleio small lifecycle_scaleio_service')
  expect(output).to include('OK')
end

def bind_service
  output = exec_command("cf bind-service #{app_name} lifecycle_scaleio_service")
  expect(output).to include('OK')

  output = exec_command("cf env #{app_name}")
  expect(output).to include('ScaleIO')
end

def unbind_service
  output = exec_command("cf unbind-service #{app_name} lifecycle_scaleio_service")
  expect(output).to include('OK')

  output = exec_command("cf env #{app_name}")
  expect(output).to_not include('ScaleIO')
end

def delete_service_instance
  output = exec_command('cf delete-service lifecycle_scaleio_service -f')
  expect(output).to include('OK')
end

def delete_not_created_service
  output = exec_command('cf delete-service notcreatedservice -f')
  puts "The output is #{output}"
  expect(output).to include('Service notcreatedservice does not exist.')
end
