{{define "ServiceGetByServiceID"}}
select * from `t_service`  where `service_id`=:ServiceID  and `deleted_at` is null;
{{end}}

{{define "ServiceGetAllByServiceIDList"}}
select * from `t_service`  where `service_id` in ({{in . .ServiceIDList}})  and `deleted_at` is null;
{{end}}


{{define "ServicePaginateWhere"}}
  
{{end}}

{{define "ServicePaginateTotal"}}
select count(*) as `count` from `t_service`  where 1=1 {{template "ServicePaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "ServicePaginate"}}
select * from `t_service`  where 1=1 {{template "ServicePaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "ServiceInsert"}}
insert into `t_service` (`service_id`,`title`,`description`,`version`,`contact_ids`,`license`,`security`,`proxy`,`variables`)values
 (:ServiceID,:Title,:Description,:Version,:ContactIds,:License,:Security,:Proxy,:Variables);
{{end}}

{{define "ServiceUpdate"}}
{{$preComma:=newPreComma}}
 update `t_service` set {{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}} 
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}} 
{{if .Version}} {{$preComma.PreComma}} `version`=:Version {{end}} 
{{if .ContactIds}} {{$preComma.PreComma}} `contact_ids`=:ContactIds {{end}} 
{{if .License}} {{$preComma.PreComma}} `license`=:License {{end}} 
{{if .Security}} {{$preComma.PreComma}} `security`=:Security {{end}} 
{{if .Proxy}} {{$preComma.PreComma}} `proxy`=:Proxy {{end}} 
{{if .Variables}} {{$preComma.PreComma}} `variables`=:Variables {{end}}  where `service_id`=:ServiceID;
{{end}}

{{define "ServiceDel"}}
update `t_service` set `deleted_at`={{currentTime .}} where `service_id`=:ServiceID;
{{end}}

{{define "APIGetByAPIID"}}
select * from `t_api`  where `api_id`=:APIID  and `deleted_at` is null;
{{end}}

{{define "APIGetAllByAPIIDList"}}
select * from `t_api`  where `api_id` in ({{in . .APIIDList}})  and `deleted_at` is null;
{{end}}


{{define "APIPaginateWhere"}}
  
{{end}}

{{define "APIPaginateTotal"}}
select count(*) as `count` from `t_api`  where 1=1 {{template "APIPaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "APIPaginate"}}
select * from `t_api`  where 1=1 {{template "APIPaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "APIInsert"}}
insert into `t_api` (`api_id`,`service_id`,`name`,`title`,`tags`,`uri`,`summary`,`description`)values
 (:APIID,:ServiceID,:Name,:Title,:Tags,:URI,:Summary,:Description);
{{end}}

{{define "APIUpdate"}}
{{$preComma:=newPreComma}}
 update `t_api` set {{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .Name}} {{$preComma.PreComma}} `name`=:Name {{end}} 
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}} 
{{if .Tags}} {{$preComma.PreComma}} `tags`=:Tags {{end}} 
{{if .URI}} {{$preComma.PreComma}} `uri`=:URI {{end}} 
{{if .Summary}} {{$preComma.PreComma}} `summary`=:Summary {{end}} 
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}}  where `api_id`=:APIID;
{{end}}

{{define "APIDel"}}
update `t_api` set `deleted_at`={{currentTime .}} where `api_id`=:APIID;
{{end}}

{{define "ParameterGetByParameterID"}}
select * from `t_parameter`  where `parameter_id`=:ParameterID  and `deleted_at` is null;
{{end}}

{{define "ParameterGetAllByParameterIDList"}}
select * from `t_parameter`  where `parameter_id` in ({{in . .ParameterIDList}})  and `deleted_at` is null;
{{end}}


{{define "ParameterPaginateWhere"}}
  
{{end}}

{{define "ParameterPaginateTotal"}}
select count(*) as `count` from `t_parameter`  where 1=1 {{template "ParameterPaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "ParameterPaginate"}}
select * from `t_parameter`  where 1=1 {{template "ParameterPaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "ParameterInsert"}}
insert into `t_parameter` (`parameter_id`,`service_id`,`api_id`,`validate_schema_id`,`full_name`,`name`,`title`,`type`,`tag`,`method`,`http_status`,`position`,`example`,`deprecated`,`required`,`serialize`,`explode`,`allow_empty_value`,`allow_reserved`,`description`)values
 (:ParameterID,:ServiceID,:APIID,:ValidateSchemaID,:FullName,:Name,:Title,:Type,:Tag,:Method,:HTTPStatus,:Position,:Example,:Deprecated,:Required,:Serialize,:Explode,:AllowEmptyValue,:AllowReserved,:Description);
{{end}}

{{define "ParameterUpdate"}}
{{$preComma:=newPreComma}}
 update `t_parameter` set {{if .ParameterID}} {{$preComma.PreComma}} `parameter_id`=:ParameterID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}} 
{{if .ValidateSchemaID}} {{$preComma.PreComma}} `validate_schema_id`=:ValidateSchemaID {{end}} 
{{if .FullName}} {{$preComma.PreComma}} `full_name`=:FullName {{end}} 
{{if .Name}} {{$preComma.PreComma}} `name`=:Name {{end}} 
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}} 
{{if .Type}} {{$preComma.PreComma}} `type`=:Type {{end}} 
{{if .Tag}} {{$preComma.PreComma}} `tag`=:Tag {{end}} 
{{if .Method}} {{$preComma.PreComma}} `method`=:Method {{end}} 
{{if .HTTPStatus}} {{$preComma.PreComma}} `http_status`=:HTTPStatus {{end}} 
{{if .Position}} {{$preComma.PreComma}} `position`=:Position {{end}} 
{{if .Example}} {{$preComma.PreComma}} `example`=:Example {{end}} 
{{if .Deprecated}} {{$preComma.PreComma}} `deprecated`=:Deprecated {{end}} 
{{if .Required}} {{$preComma.PreComma}} `required`=:Required {{end}} 
{{if .Serialize}} {{$preComma.PreComma}} `serialize`=:Serialize {{end}} 
{{if .Explode}} {{$preComma.PreComma}} `explode`=:Explode {{end}} 
{{if .AllowEmptyValue}} {{$preComma.PreComma}} `allow_empty_value`=:AllowEmptyValue {{end}} 
{{if .AllowReserved}} {{$preComma.PreComma}} `allow_reserved`=:AllowReserved {{end}} 
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}}  where `parameter_id`=:ParameterID;
{{end}}

{{define "ParameterDel"}}
update `t_parameter` set `deleted_at`={{currentTime .}} where `parameter_id`=:ParameterID;
{{end}}

{{define "ValidateSchemaGetByValidateSchemaID"}}
select * from `t_validate_schema`  where `validate_schema_id`=:ValidateSchemaID  and `deleted_at` is null;
{{end}}

{{define "ValidateSchemaGetAllByValidateSchemaIDList"}}
select * from `t_validate_schema`  where `validate_schema_id` in ({{in . .ValidateSchemaIDList}})  and `deleted_at` is null;
{{end}}


{{define "ValidateSchemaPaginateWhere"}}
  
{{end}}

{{define "ValidateSchemaPaginateTotal"}}
select count(*) as `count` from `t_validate_schema`  where 1=1 {{template "ValidateSchemaPaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "ValidateSchemaPaginate"}}
select * from `t_validate_schema`  where 1=1 {{template "ValidateSchemaPaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "ValidateSchemaInsert"}}
insert into `t_validate_schema` (`validate_schema_id`,`service_id`,`description`,`remark`,`type`,`example`,`deprecated`,`required`,`enum`,`enum_names`,`enum_titles`,`format`,`default`,`nullable`,`multiple_of`,`maxnum`,`exclusive_maximum`,`minimum`,`exclusive_minimum`,`max_length`,`min_length`,`pattern`,`max_items`,`min_items`,`unique_items`,`max_properties`,`min_properties`,`all_of`,`one_of`,`any_of`,`allow_empty_value`,`allow_reserved`,`not`,`additional_properties`,`discriminator`,`read_only`,`write_only`,`xml`,`external_docs`,`external_pros`,`extensions`,`summary`)values
 (:ValidateSchemaID,:ServiceID,:Description,:Remark,:Type,:Example,:Deprecated,:Required,:Enum,:EnumNames,:EnumTitles,:Format,:Default,:Nullable,:MultipleOf,:Maxnum,:ExclusiveMaximum,:Minimum,:ExclusiveMinimum,:MaxLength,:MinLength,:Pattern,:MaxItems,:MinItems,:UniqueItems,:MaxProperties,:MinProperties,:AllOf,:OneOf,:AnyOf,:AllowEmptyValue,:AllowReserved,:Not,:AdditionalProperties,:Discriminator,:ReadOnly,:WriteOnly,:XML,:ExternalDocs,:ExternalPros,:Extensions,:Summary);
{{end}}

{{define "ValidateSchemaUpdate"}}
{{$preComma:=newPreComma}}
 update `t_validate_schema` set {{if .ValidateSchemaID}} {{$preComma.PreComma}} `validate_schema_id`=:ValidateSchemaID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}} 
{{if .Remark}} {{$preComma.PreComma}} `remark`=:Remark {{end}} 
{{if .Type}} {{$preComma.PreComma}} `type`=:Type {{end}} 
{{if .Example}} {{$preComma.PreComma}} `example`=:Example {{end}} 
{{if .Deprecated}} {{$preComma.PreComma}} `deprecated`=:Deprecated {{end}} 
{{if .Required}} {{$preComma.PreComma}} `required`=:Required {{end}} 
{{if .Enum}} {{$preComma.PreComma}} `enum`=:Enum {{end}} 
{{if .EnumNames}} {{$preComma.PreComma}} `enum_names`=:EnumNames {{end}} 
{{if .EnumTitles}} {{$preComma.PreComma}} `enum_titles`=:EnumTitles {{end}} 
{{if .Format}} {{$preComma.PreComma}} `format`=:Format {{end}} 
{{if .Default}} {{$preComma.PreComma}} `default`=:Default {{end}} 
{{if .Nullable}} {{$preComma.PreComma}} `nullable`=:Nullable {{end}} 
{{if .MultipleOf}} {{$preComma.PreComma}} `multiple_of`=:MultipleOf {{end}} 
{{if .Maxnum}} {{$preComma.PreComma}} `maxnum`=:Maxnum {{end}} 
{{if .ExclusiveMaximum}} {{$preComma.PreComma}} `exclusive_maximum`=:ExclusiveMaximum {{end}} 
{{if .Minimum}} {{$preComma.PreComma}} `minimum`=:Minimum {{end}} 
{{if .ExclusiveMinimum}} {{$preComma.PreComma}} `exclusive_minimum`=:ExclusiveMinimum {{end}} 
{{if .MaxLength}} {{$preComma.PreComma}} `max_length`=:MaxLength {{end}} 
{{if .MinLength}} {{$preComma.PreComma}} `min_length`=:MinLength {{end}} 
{{if .Pattern}} {{$preComma.PreComma}} `pattern`=:Pattern {{end}} 
{{if .MaxItems}} {{$preComma.PreComma}} `max_items`=:MaxItems {{end}} 
{{if .MinItems}} {{$preComma.PreComma}} `min_items`=:MinItems {{end}} 
{{if .UniqueItems}} {{$preComma.PreComma}} `unique_items`=:UniqueItems {{end}} 
{{if .MaxProperties}} {{$preComma.PreComma}} `max_properties`=:MaxProperties {{end}} 
{{if .MinProperties}} {{$preComma.PreComma}} `min_properties`=:MinProperties {{end}} 
{{if .AllOf}} {{$preComma.PreComma}} `all_of`=:AllOf {{end}} 
{{if .OneOf}} {{$preComma.PreComma}} `one_of`=:OneOf {{end}} 
{{if .AnyOf}} {{$preComma.PreComma}} `any_of`=:AnyOf {{end}} 
{{if .AllowEmptyValue}} {{$preComma.PreComma}} `allow_empty_value`=:AllowEmptyValue {{end}} 
{{if .AllowReserved}} {{$preComma.PreComma}} `allow_reserved`=:AllowReserved {{end}} 
{{if .Not}} {{$preComma.PreComma}} `not`=:Not {{end}} 
{{if .AdditionalProperties}} {{$preComma.PreComma}} `additional_properties`=:AdditionalProperties {{end}} 
{{if .Discriminator}} {{$preComma.PreComma}} `discriminator`=:Discriminator {{end}} 
{{if .ReadOnly}} {{$preComma.PreComma}} `read_only`=:ReadOnly {{end}} 
{{if .WriteOnly}} {{$preComma.PreComma}} `write_only`=:WriteOnly {{end}} 
{{if .XML}} {{$preComma.PreComma}} `xml`=:XML {{end}} 
{{if .ExternalDocs}} {{$preComma.PreComma}} `external_docs`=:ExternalDocs {{end}} 
{{if .ExternalPros}} {{$preComma.PreComma}} `external_pros`=:ExternalPros {{end}} 
{{if .Extensions}} {{$preComma.PreComma}} `extensions`=:Extensions {{end}} 
{{if .Summary}} {{$preComma.PreComma}} `summary`=:Summary {{end}}  where `validate_schema_id`=:ValidateSchemaID;
{{end}}

{{define "ValidateSchemaDel"}}
update `t_validate_schema` set `deleted_at`={{currentTime .}} where `validate_schema_id`=:ValidateSchemaID;
{{end}}

{{define "ExampleGetByExampleID"}}
select * from `t_example`  where `example_id`=:ExampleID  and `deleted_at` is null;
{{end}}

{{define "ExampleGetAllByExampleIDList"}}
select * from `t_example`  where `example_id` in ({{in . .ExampleIDList}})  and `deleted_at` is null;
{{end}}


{{define "ExamplePaginateWhere"}}
  
{{end}}

{{define "ExamplePaginateTotal"}}
select count(*) as `count` from `t_example`  where 1=1 {{template "ExamplePaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "ExamplePaginate"}}
select * from `t_example`  where 1=1 {{template "ExamplePaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "ExampleInsert"}}
insert into `t_example` (`example_id`,`service_id`,`api_id`,`tag`,`title`,`summary`,`request`,`response`)values
 (:ExampleID,:ServiceID,:APIID,:Tag,:Title,:Summary,:Request,:Response);
{{end}}

{{define "ExampleUpdate"}}
{{$preComma:=newPreComma}}
 update `t_example` set {{if .ExampleID}} {{$preComma.PreComma}} `example_id`=:ExampleID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}} 
{{if .Tag}} {{$preComma.PreComma}} `tag`=:Tag {{end}} 
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}} 
{{if .Summary}} {{$preComma.PreComma}} `summary`=:Summary {{end}} 
{{if .Request}} {{$preComma.PreComma}} `request`=:Request {{end}} 
{{if .Response}} {{$preComma.PreComma}} `response`=:Response {{end}}  where `example_id`=:ExampleID;
{{end}}

{{define "ExampleDel"}}
update `t_example` set `deleted_at`={{currentTime .}} where `example_id`=:ExampleID;
{{end}}

{{define "MarkdownGetByMarkdownID"}}
select * from `t_markdown`  where `markdown_id`=:MarkdownID  and `deleted_at` is null;
{{end}}

{{define "MarkdownGetAllByMarkdownIDList"}}
select * from `t_markdown`  where `markdown_id` in ({{in . .MarkdownIDList}})  and `deleted_at` is null;
{{end}}


{{define "MarkdownPaginateWhere"}}
  
{{end}}

{{define "MarkdownPaginateTotal"}}
select count(*) as `count` from `t_markdown`  where 1=1 {{template "MarkdownPaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "MarkdownPaginate"}}
select * from `t_markdown`  where 1=1 {{template "MarkdownPaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "MarkdownInsert"}}
insert into `t_markdown` (`markdown_id`,`service_id`,`api_id`,`name`,`title`,`markdown`,`content`,`owner_id`,`owner_name`)values
 (:MarkdownID,:ServiceID,:APIID,:Name,:Title,:Markdown,:Content,:OwnerID,:OwnerName);
{{end}}

{{define "MarkdownUpdate"}}
{{$preComma:=newPreComma}}
 update `t_markdown` set {{if .MarkdownID}} {{$preComma.PreComma}} `markdown_id`=:MarkdownID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}} 
{{if .Name}} {{$preComma.PreComma}} `name`=:Name {{end}} 
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}} 
{{if .Markdown}} {{$preComma.PreComma}} `markdown`=:Markdown {{end}} 
{{if .Content}} {{$preComma.PreComma}} `content`=:Content {{end}} 
{{if .OwnerID}} {{$preComma.PreComma}} `owner_id`=:OwnerID {{end}} 
{{if .OwnerName}} {{$preComma.PreComma}} `owner_name`=:OwnerName {{end}}  where `markdown_id`=:MarkdownID;
{{end}}

{{define "MarkdownDel"}}
update `t_markdown` set `deleted_at`={{currentTime .}} where `markdown_id`=:MarkdownID;
{{end}}

{{define "ServerGetByServerID"}}
select * from `t_server`  where `server_id`=:ServerID  and `deleted_at` is null;
{{end}}

{{define "ServerGetAllByServerIDList"}}
select * from `t_server`  where `server_id` in ({{in . .ServerIDList}})  and `deleted_at` is null;
{{end}}


{{define "ServerPaginateWhere"}}
  
{{end}}

{{define "ServerPaginateTotal"}}
select count(*) as `count` from `t_server`  where 1=1 {{template "ServerPaginateWhere" .}}   and `deleted_at` is null;
{{end}}

  {{define "ServerPaginate"}}
select * from `t_server`  where 1=1 {{template "ServerPaginateWhere" .}}   and `deleted_at` is null order by `updated_at` desc  limit :Offset,:Limit ;
{{end}}


{{define "ServerInsert"}}
insert into `t_server` (`server_id`,`service_id`,`url`,`description`,`proxy`,`extension_ids`)values
 (:ServerID,:ServiceID,:URL,:Description,:Proxy,:ExtensionIds);
{{end}}

{{define "ServerUpdate"}}
{{$preComma:=newPreComma}}
 update `t_server` set {{if .ServerID}} {{$preComma.PreComma}} `server_id`=:ServerID {{end}} 
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}} 
{{if .URL}} {{$preComma.PreComma}} `url`=:URL {{end}} 
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}} 
{{if .Proxy}} {{$preComma.PreComma}} `proxy`=:Proxy {{end}} 
{{if .ExtensionIds}} {{$preComma.PreComma}} `extension_ids`=:ExtensionIds {{end}}  where `server_id`=:ServerID;
{{end}}

{{define "ServerDel"}}
update `t_server` set `deleted_at`={{currentTime .}} where `server_id`=:ServerID;
{{end}}
