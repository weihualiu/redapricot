-module(agent_pack).

-export([sendLog/1, sendLog/2,
	binary_currdate/0
]).


sendLog([AppName, Tag, ExecTime, UseTime, Content, {Host, Port}]) ->
	try
		sendLog([AppName, Tag, [], ExecTime, UseTime, Content], [agent_socket, send, {Host, Port}])
	catch
		ErrT:ErrM ->
			io:format("agent error type:~p, message:~p~n", [ErrT, ErrM])
	end,
	ok;
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content, {Host, Port}]) ->
	try
		sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [agent_socket, send, {Host, Port}])
	catch
		ErrT:ErrM ->
			io:format("agent error type:~p, message:~p~n", [ErrT, ErrM])
	end,
	ok.


%% 上送日志耗时统计数据
%% M,F : atom
%% 直接返回M:F(A,B)的结果
%% example:  agent_pack:sendLog(["muip", "testtestest_MAIL0001_1", agent_pack:binary_currdate(), 120, "ceshi jiekou"], [agent_socket, send, {"182.207.176.150", 8090}]).
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) when is_list(AppName) ->
	sendLog([list_to_binary(AppName), Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]);
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) when is_list(Tag) ->
	sendLog([AppName, list_to_binary(Tag), ClientTag, ExecTime, UseTime, Content], [M,F,A]);
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) when is_list(ClientTag) ->
	sendLog([AppName, list_to_binary(ClientTag), ClientTag, ExecTime, UseTime, Content], [M,F,A]);
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) when is_list(ExecTime) ->
	sendLog([AppName, Tag, ClientTag, list_to_binary(ExecTime), UseTime, Content], [M,F,A]);
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) when is_list(Content) ->
	sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, list_to_binary(Content)], [M,F,A]);
sendLog([AppName, Tag, ClientTag, ExecTime, UseTime, Content], [M,F,A]) ->
	%% 数据组装
	%%  UseTime 4bytes
	%%  Tag  100bytes
	%%  AppName 20bytes
	%%  ExecTime 7bytes
	%%  Content Nbytes
	TagB = binary_addzero(Tag, 100),
	ClientB = binary_addzero(ClientTag, 100),
	AppNameB = binary_addzero(AppName, 20),
	ExecTimeB = binary_addzero(ExecTime, 7),
	Context = <<UseTime:32/big, TagB/binary, ClientB/binary, AppNameB/binary, ExecTimeB/binary, Content/binary>>,
	Bin = package_build(<<16#00>>, Context),
	M:F(A, Bin).


binary_addzero(Bin, Num) when size(Bin) == Num ->
    Bin;
binary_addzero(Bin, Num) when size(Bin) > Num ->
    binary:part(Bin,0, Num);
binary_addzero(Bin, Num) ->
    binary_addzero_t(Bin,Num-size(Bin)).

binary_addzero_t(Bin, 0) ->
    Bin;
binary_addzero_t(Bin, Num) ->
    binary_addzero_t(<<Bin/binary, 16#0>>, Num - 1).

binary_currdate() ->
    %% binary 7个字节存储 yyyymmddhhMMSS
    {{Year, Month, Day}, {Hour, Min, Sec}} = calendar:now_to_local_time(now()),
    YearT = Year rem 100,
    YearH = trunc(Year/100),
    <<YearH:8, YearT:8, Month:8, Day:8, Hour:8, Min:8, Sec:8 >>.

package_build(Protocol, Context) ->
	Len = 1 + 4 + 1 + 7 + 1 + size(Context),
	Currdate = binary_currdate(),
	Bin = <<16#F0, Len:32/big, Protocol/binary, Currdate/binary, Context/binary, 16#FE>>,
	Bin.
