-module(agent_socket).

-define(TCP_BUFFER_SIZE, 10240).
-define(SENDTIMEOUT, 36000).
-define(REQUEST_ERROR, "agent socket failed!").
-define(REQUEST_TIMEOUT, "agent socket timeout!").

-export([send/2]).

send({Host, Port}, Data) ->
	TCPOptions = [binary, {packet, raw}, {active, false}, {reuseaddr, true},        % 组装tcp选项, 
		{buffer, ?TCP_BUFFER_SIZE}, {send_timeout, ?SENDTIMEOUT}],
	case gen_tcp:connect(Host, Port, TCPOptions, ?SENDTIMEOUT) of                    % 建立连接
		{ok, Socket} ->                                                                 % 建立连接成功
			ResData = send_and_recv(Socket, Data, ?SENDTIMEOUT),                         % 发收消息
			gen_tcp:close(Socket),                                                      % 关闭连接
			ResData;
		{error, Reason} ->                                                              % 建立连接失败
			io:format("reason:~p~n", [Reason]),
			throw(?REQUEST_ERROR)
	end.

send_and_recv(Socket, ReqData, Timeout) ->
	case gen_tcp:send(Socket, ReqData) of
		ok ->
			case do_recv(Socket, Timeout, []) of
				{ok, ResData} ->
					ResData;
				{error, timeout} ->
					gen_tcp:close(Socket),
					io:format("request bank system timeout", []),
					throw(?REQUEST_TIMEOUT);
				{error, Reason} ->
					gen_tcp:close(Socket),
					io:format("request bank system faild:~p", [Reason]),
					throw(?REQUEST_ERROR)
			end;
		{error, Reason} ->
			gen_tcp:close(Socket),
			io:format("send to bank system faild:~p", [Reason]),
			throw(?REQUEST_ERROR)
	end.

%% @doc 发收消息
%% @end
do_recv(Sock, Timeout, Bs) ->
	case gen_tcp:recv(Sock, 0, Timeout) of
		{ok, B} ->
			do_recv(Sock, Timeout, [B|Bs]);
		{error, closed} ->
			{ok, list_to_binary(lists:reverse(Bs))};
		Error ->
			io:format("Error = ~p~n", [Error]),
			throw(?REQUEST_ERROR)
	end.
