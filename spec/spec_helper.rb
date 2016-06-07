# spec/spec_helper.rb
require 'rspec'
require 'json'

def endpoint
  ENV['CF_ENDPOINT']
end

def build_pack
  "go_buildpack"
end

def project_path
  File.expand_path('../..', __FILE__)
end
