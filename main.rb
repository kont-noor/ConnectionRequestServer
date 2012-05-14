require 'sinatra'

post '/request_permission_to_connect' do
  request_permission_to_connect params
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

def request_permission_to_connect param
  activation_code = param[:activation_code]
  device_id = param[:device_id]
  client_version = param[:client_version]
  os_version = param[:os_version]

  response = {
    1 =>'Approved',
    400 => 'Sorry, your account is currently connected from another computer. You can use our service from multiple computers, but each account can only be connected to our network from one computer at a time. To connect from this computer now, please buy an additional account.',
    401 => "Missing parameters. Sorry, we've made a note to fix this. Please try again and contact support if you continue to see this error.",
    500 => 'Sorry, unknown error. Please try again and contact support if you continue to see this error.'
  }
  
  if (activation_code.nil? or device_id.nil?)
    response_code = 401
  else
    response_code = 500
  end

  "<connection_request_response><code>#{response_code}</code><message>#{response[response_code]}</message></connection_request_response>"
end

def disconnect param
end

def heartbeat param
end
