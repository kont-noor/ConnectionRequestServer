require 'sinatra'
require 'yaml'

class ConnectionRequestServer < Sinatra::Base

  post '/request_permission_to_connect' do
    response_code = request_permission_to_connect params
    content_type 'text/xml'

    message =  get_messages
    erb :connection_response, :locals => {:code => response_code, :message => message[response_code]}
  end

  post '/disconnect' do
    disconnect params
    'ok'
  end

  post '/heartbeat' do
    heartbeat params
    'ok'
  end

  error 400..510 do
    '404'
  end

  def get_messages
    messages_yaml = File.read(File.dirname(__FILE__) + '/config/messages.yml')
    YAML.load(messages_yaml)
  end

  def request_permission_to_connect param
    activation_code = param[:activation_code]
    device_id = param[:device_id]
    client_version = param[:client_version]
    os_version = param[:os_version]

    if (activation_code.nil? or device_id.nil?)
      return 401
    else
      return 500
    end

  end

  def disconnect param
  end

  def heartbeat param
  end
end
