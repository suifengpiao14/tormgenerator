package example

import "github.com/suifengpiao14/torm"

type ServiceGetByServiceIDEntity struct {
	ServiceID string //服务标识

	torm.VolumeMap
}

func (t *ServiceGetByServiceIDEntity) TplName() string {
	return "ServiceGetByServiceID"
}

func (t *ServiceGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type ServiceGetAllByServiceIDListEntity struct {
	ServiceIDList []string //

	torm.VolumeMap
}

func (t *ServiceGetAllByServiceIDListEntity) TplName() string {
	return "ServiceGetAllByServiceIDList"
}

func (t *ServiceGetAllByServiceIDListEntity) TplType() string {
	return "sql_select"
}

type ServicePaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *ServicePaginateWhereEntity) TplName() string {
	return "ServicePaginateWhere"
}

func (t *ServicePaginateWhereEntity) TplType() string {
	return "text"
}

type ServicePaginateTotalEntity struct {
	ServicePaginateWhereEntity //

	torm.VolumeMap
}

func (t *ServicePaginateTotalEntity) TplName() string {
	return "ServicePaginateTotal"
}

func (t *ServicePaginateTotalEntity) TplType() string {
	return "sql_select"
}

type ServicePaginateEntity struct {
	Limit int //

	Offset int //

	ServicePaginateWhereEntity //

	torm.VolumeMap
}

func (t *ServicePaginateEntity) TplName() string {
	return "ServicePaginate"
}

func (t *ServicePaginateEntity) TplType() string {
	return "sql_select"
}

type ServiceInsertEntity struct {
	ContactIds string //联系人

	Description string //简介

	License string //协议

	Proxy string //代理地址

	Security string //鉴权

	ServiceID string //服务标识

	Title string //名称(冗余local.zh)

	Variables string //json字符串

	Version string //版本

	torm.VolumeMap
}

func (t *ServiceInsertEntity) TplName() string {
	return "ServiceInsert"
}

func (t *ServiceInsertEntity) TplType() string {
	return "sql_insert"
}

type ServiceUpdateEntity struct {
	ContactIds string //联系人

	Description string //简介

	License string //协议

	Proxy string //代理地址

	Security string //鉴权

	ServiceID string //服务标识

	Title string //名称(冗余local.zh)

	Variables string //json字符串

	Version string //版本

	torm.VolumeMap
}

func (t *ServiceUpdateEntity) TplName() string {
	return "ServiceUpdate"
}

func (t *ServiceUpdateEntity) TplType() string {
	return "sql_update"
}

type ServiceDelEntity struct {
	ServiceID string //服务标识

	torm.VolumeMap
}

func (t *ServiceDelEntity) TplName() string {
	return "ServiceDel"
}

func (t *ServiceDelEntity) TplType() string {
	return "sql_update"
}

type APIGetByAPIIDEntity struct {
	APIID string //api 文档标识

	torm.VolumeMap
}

func (t *APIGetByAPIIDEntity) TplName() string {
	return "APIGetByAPIID"
}

func (t *APIGetByAPIIDEntity) TplType() string {
	return "sql_select"
}

type APIGetAllByAPIIDListEntity struct {
	APIIDList []string //

	torm.VolumeMap
}

func (t *APIGetAllByAPIIDListEntity) TplName() string {
	return "APIGetAllByAPIIDList"
}

func (t *APIGetAllByAPIIDListEntity) TplType() string {
	return "sql_select"
}

type APIPaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *APIPaginateWhereEntity) TplName() string {
	return "APIPaginateWhere"
}

func (t *APIPaginateWhereEntity) TplType() string {
	return "text"
}

type APIPaginateTotalEntity struct {
	APIPaginateWhereEntity //

	torm.VolumeMap
}

func (t *APIPaginateTotalEntity) TplName() string {
	return "APIPaginateTotal"
}

func (t *APIPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type APIPaginateEntity struct {
	APIPaginateWhereEntity //

	Limit int //

	Offset int //

	torm.VolumeMap
}

func (t *APIPaginateEntity) TplName() string {
	return "APIPaginate"
}

func (t *APIPaginateEntity) TplType() string {
	return "sql_select"
}

type APIInsertEntity struct {
	APIID string //api 文档标识

	Description string //简介

	Name string //名称(冗余local.en)

	ServiceID string //服务标识

	Summary string //摘要

	Tags string //标签ID集合

	Title string //名称(冗余local.zh)

	URI string //路径

	torm.VolumeMap
}

func (t *APIInsertEntity) TplName() string {
	return "APIInsert"
}

func (t *APIInsertEntity) TplType() string {
	return "sql_insert"
}

type APIUpdateEntity struct {
	APIID string //api 文档标识

	Description string //简介

	Name string //名称(冗余local.en)

	ServiceID string //服务标识

	Summary string //摘要

	Tags string //标签ID集合

	Title string //名称(冗余local.zh)

	URI string //路径

	torm.VolumeMap
}

func (t *APIUpdateEntity) TplName() string {
	return "APIUpdate"
}

func (t *APIUpdateEntity) TplType() string {
	return "sql_update"
}

type APIDelEntity struct {
	APIID string //api 文档标识

	torm.VolumeMap
}

func (t *APIDelEntity) TplName() string {
	return "APIDel"
}

func (t *APIDelEntity) TplType() string {
	return "sql_update"
}

type ParameterGetByParameterIDEntity struct {
	ParameterID string //参数标识

	torm.VolumeMap
}

func (t *ParameterGetByParameterIDEntity) TplName() string {
	return "ParameterGetByParameterID"
}

func (t *ParameterGetByParameterIDEntity) TplType() string {
	return "sql_select"
}

type ParameterGetAllByParameterIDListEntity struct {
	ParameterIDList []string //

	torm.VolumeMap
}

func (t *ParameterGetAllByParameterIDListEntity) TplName() string {
	return "ParameterGetAllByParameterIDList"
}

func (t *ParameterGetAllByParameterIDListEntity) TplType() string {
	return "sql_select"
}

type ParameterPaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *ParameterPaginateWhereEntity) TplName() string {
	return "ParameterPaginateWhere"
}

func (t *ParameterPaginateWhereEntity) TplType() string {
	return "text"
}

type ParameterPaginateTotalEntity struct {
	ParameterPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ParameterPaginateTotalEntity) TplName() string {
	return "ParameterPaginateTotal"
}

func (t *ParameterPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type ParameterPaginateEntity struct {
	Limit int //

	Offset int //

	ParameterPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ParameterPaginateEntity) TplName() string {
	return "ParameterPaginate"
}

func (t *ParameterPaginateEntity) TplType() string {
	return "sql_select"
}

type ParameterInsertEntity struct {
	APIID string //api 文档标识

	AllowEmptyValue string //是否容许空值(true-是,false-否)

	AllowReserved string //特殊字符是否容许出现在uri参数中(true-是,false-否)

	Deprecated string //是否弃用(true-是,false-否)

	Description string //简介

	Example string //案例

	Explode string //对象的key,是否单独成参数方式,参照openapi parameters.explode(true-是,false-否)

	FullName string //全称

	HTTPStatus string //http状态码

	Method string //请求方法(get-GET,post-POST,put-PUT,delete-DELETE,head-HEAD)

	Name string //名称(冗余local.en)

	ParameterID string //参数标识

	Position string //参数所在的位置(body-body,head-head,path-path,query-query,cookie-cookie)

	Required string //是否必须(true-是,false-否)

	Serialize string //对数组、对象序列化方法,参照openapi parameters.style

	ServiceID string //服务标识

	Tag string //所属标签

	Title string //名称(冗余local.zh)

	Type string //参数类型:(string-字符串,int-整型,number-数字,array-数组,object-对象)

	ValidateSchemaID string //验证规则标识

	torm.VolumeMap
}

func (t *ParameterInsertEntity) TplName() string {
	return "ParameterInsert"
}

func (t *ParameterInsertEntity) TplType() string {
	return "sql_insert"
}

type ParameterUpdateEntity struct {
	APIID string //api 文档标识

	AllowEmptyValue string //是否容许空值(true-是,false-否)

	AllowReserved string //特殊字符是否容许出现在uri参数中(true-是,false-否)

	Deprecated string //是否弃用(true-是,false-否)

	Description string //简介

	Example string //案例

	Explode string //对象的key,是否单独成参数方式,参照openapi parameters.explode(true-是,false-否)

	FullName string //全称

	HTTPStatus string //http状态码

	Method string //请求方法(get-GET,post-POST,put-PUT,delete-DELETE,head-HEAD)

	Name string //名称(冗余local.en)

	ParameterID string //参数标识

	Position string //参数所在的位置(body-body,head-head,path-path,query-query,cookie-cookie)

	Required string //是否必须(true-是,false-否)

	Serialize string //对数组、对象序列化方法,参照openapi parameters.style

	ServiceID string //服务标识

	Tag string //所属标签

	Title string //名称(冗余local.zh)

	Type string //参数类型:(string-字符串,int-整型,number-数字,array-数组,object-对象)

	ValidateSchemaID string //验证规则标识

	torm.VolumeMap
}

func (t *ParameterUpdateEntity) TplName() string {
	return "ParameterUpdate"
}

func (t *ParameterUpdateEntity) TplType() string {
	return "sql_update"
}

type ParameterDelEntity struct {
	ParameterID string //参数标识

	torm.VolumeMap
}

func (t *ParameterDelEntity) TplName() string {
	return "ParameterDel"
}

func (t *ParameterDelEntity) TplType() string {
	return "sql_update"
}

type ValidateSchemaGetByValidateSchemaIDEntity struct {
	ValidateSchemaID string //验证规则标识

	torm.VolumeMap
}

func (t *ValidateSchemaGetByValidateSchemaIDEntity) TplName() string {
	return "ValidateSchemaGetByValidateSchemaID"
}

func (t *ValidateSchemaGetByValidateSchemaIDEntity) TplType() string {
	return "sql_select"
}

type ValidateSchemaGetAllByValidateSchemaIDListEntity struct {
	ValidateSchemaIDList []string //

	torm.VolumeMap
}

func (t *ValidateSchemaGetAllByValidateSchemaIDListEntity) TplName() string {
	return "ValidateSchemaGetAllByValidateSchemaIDList"
}

func (t *ValidateSchemaGetAllByValidateSchemaIDListEntity) TplType() string {
	return "sql_select"
}

type ValidateSchemaPaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *ValidateSchemaPaginateWhereEntity) TplName() string {
	return "ValidateSchemaPaginateWhere"
}

func (t *ValidateSchemaPaginateWhereEntity) TplType() string {
	return "text"
}

type ValidateSchemaPaginateTotalEntity struct {
	ValidateSchemaPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ValidateSchemaPaginateTotalEntity) TplName() string {
	return "ValidateSchemaPaginateTotal"
}

func (t *ValidateSchemaPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type ValidateSchemaPaginateEntity struct {
	Limit int //

	Offset int //

	ValidateSchemaPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ValidateSchemaPaginateEntity) TplName() string {
	return "ValidateSchemaPaginate"
}

func (t *ValidateSchemaPaginateEntity) TplType() string {
	return "sql_select"
}

type ValidateSchemaInsertEntity struct {
	AdditionalProperties string //boolean

	AllOf string //所有

	AllowEmptyValue string //是否容许空值(true-是,false-否)

	AllowReserved string //特殊字符是否容许出现在uri参数中(true-是,false-否)

	AnyOf string //任何一个SchemaID

	Default string //默认值

	Deprecated string //是否弃用(true-是,false-否)

	Description string //简介

	Discriminator string //schema鉴别

	Enum string //枚举值

	EnumNames string //枚举名称

	EnumTitles string //枚举标题

	Example string //案例

	ExclusiveMaximum string //是否不包含最大项(true-是,false-否)

	ExclusiveMinimum string //是否不包含最小项(true-是,false-否)

	Extensions string //扩展字段

	ExternalDocs string //附加文档

	ExternalPros string //附加文档

	Format string //格式

	MaxItems int //最大项数

	MaxLength int //最大长度

	MaxProperties int //最多属性项

	Maxnum int //最大值

	MinItems int //最小项数

	MinLength int //最小长度

	MinProperties int //最少属性项

	Minimum int //最小值

	MultipleOf int //倍数

	Not string //不包含的schemaID

	Nullable string //是否可以为空(true-是,false-否)

	OneOf string //只满足一个

	Pattern string //正则表达式

	ReadOnly string //是否只读(true-是,false-否)

	Remark string //备注

	Required string //是否必须(true-是,false-否)

	ServiceID string //服务标识

	Summary string //摘要

	Type string //参数类型:(string-字符串,int-整型,number-数字,array-数组,object-对象)

	UniqueItems string //所有项是否需要唯一(true-是,false-否)

	ValidateSchemaID string //验证规则标识

	WriteOnly string //是否只写(true-是,false-否)

	XML string //xml对象

	torm.VolumeMap
}

func (t *ValidateSchemaInsertEntity) TplName() string {
	return "ValidateSchemaInsert"
}

func (t *ValidateSchemaInsertEntity) TplType() string {
	return "sql_insert"
}

type ValidateSchemaUpdateEntity struct {
	AdditionalProperties string //boolean

	AllOf string //所有

	AllowEmptyValue string //是否容许空值(true-是,false-否)

	AllowReserved string //特殊字符是否容许出现在uri参数中(true-是,false-否)

	AnyOf string //任何一个SchemaID

	Default string //默认值

	Deprecated string //是否弃用(true-是,false-否)

	Description string //简介

	Discriminator string //schema鉴别

	Enum string //枚举值

	EnumNames string //枚举名称

	EnumTitles string //枚举标题

	Example string //案例

	ExclusiveMaximum string //是否不包含最大项(true-是,false-否)

	ExclusiveMinimum string //是否不包含最小项(true-是,false-否)

	Extensions string //扩展字段

	ExternalDocs string //附加文档

	ExternalPros string //附加文档

	Format string //格式

	MaxItems int //最大项数

	MaxLength int //最大长度

	MaxProperties int //最多属性项

	Maxnum int //最大值

	MinItems int //最小项数

	MinLength int //最小长度

	MinProperties int //最少属性项

	Minimum int //最小值

	MultipleOf int //倍数

	Not string //不包含的schemaID

	Nullable string //是否可以为空(true-是,false-否)

	OneOf string //只满足一个

	Pattern string //正则表达式

	ReadOnly string //是否只读(true-是,false-否)

	Remark string //备注

	Required string //是否必须(true-是,false-否)

	ServiceID string //服务标识

	Summary string //摘要

	Type string //参数类型:(string-字符串,int-整型,number-数字,array-数组,object-对象)

	UniqueItems string //所有项是否需要唯一(true-是,false-否)

	ValidateSchemaID string //验证规则标识

	WriteOnly string //是否只写(true-是,false-否)

	XML string //xml对象

	torm.VolumeMap
}

func (t *ValidateSchemaUpdateEntity) TplName() string {
	return "ValidateSchemaUpdate"
}

func (t *ValidateSchemaUpdateEntity) TplType() string {
	return "sql_update"
}

type ValidateSchemaDelEntity struct {
	ValidateSchemaID string //验证规则标识

	torm.VolumeMap
}

func (t *ValidateSchemaDelEntity) TplName() string {
	return "ValidateSchemaDel"
}

func (t *ValidateSchemaDelEntity) TplType() string {
	return "sql_update"
}

type ExampleGetByExampleIDEntity struct {
	ExampleID string //测试用例标识

	torm.VolumeMap
}

func (t *ExampleGetByExampleIDEntity) TplName() string {
	return "ExampleGetByExampleID"
}

func (t *ExampleGetByExampleIDEntity) TplType() string {
	return "sql_select"
}

type ExampleGetAllByExampleIDListEntity struct {
	ExampleIDList []string //

	torm.VolumeMap
}

func (t *ExampleGetAllByExampleIDListEntity) TplName() string {
	return "ExampleGetAllByExampleIDList"
}

func (t *ExampleGetAllByExampleIDListEntity) TplType() string {
	return "sql_select"
}

type ExamplePaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *ExamplePaginateWhereEntity) TplName() string {
	return "ExamplePaginateWhere"
}

func (t *ExamplePaginateWhereEntity) TplType() string {
	return "text"
}

type ExamplePaginateTotalEntity struct {
	ExamplePaginateWhereEntity //

	torm.VolumeMap
}

func (t *ExamplePaginateTotalEntity) TplName() string {
	return "ExamplePaginateTotal"
}

func (t *ExamplePaginateTotalEntity) TplType() string {
	return "sql_select"
}

type ExamplePaginateEntity struct {
	ExamplePaginateWhereEntity //

	Limit int //

	Offset int //

	torm.VolumeMap
}

func (t *ExamplePaginateEntity) TplName() string {
	return "ExamplePaginate"
}

func (t *ExamplePaginateEntity) TplType() string {
	return "sql_select"
}

type ExampleInsertEntity struct {
	APIID string //api 文档标识

	ExampleID string //测试用例标识

	Request string //

	Response string //

	ServiceID string //服务标识

	Summary string //摘要

	Tag string //所属标签

	Title string //名称(冗余local.zh)

	torm.VolumeMap
}

func (t *ExampleInsertEntity) TplName() string {
	return "ExampleInsert"
}

func (t *ExampleInsertEntity) TplType() string {
	return "sql_insert"
}

type ExampleUpdateEntity struct {
	APIID string //api 文档标识

	ExampleID string //测试用例标识

	Request string //

	Response string //

	ServiceID string //服务标识

	Summary string //摘要

	Tag string //所属标签

	Title string //名称(冗余local.zh)

	torm.VolumeMap
}

func (t *ExampleUpdateEntity) TplName() string {
	return "ExampleUpdate"
}

func (t *ExampleUpdateEntity) TplType() string {
	return "sql_update"
}

type ExampleDelEntity struct {
	ExampleID string //测试用例标识

	torm.VolumeMap
}

func (t *ExampleDelEntity) TplName() string {
	return "ExampleDel"
}

func (t *ExampleDelEntity) TplType() string {
	return "sql_update"
}

type MarkdownGetByMarkdownIDEntity struct {
	MarkdownID string //markdown 文档标识

	torm.VolumeMap
}

func (t *MarkdownGetByMarkdownIDEntity) TplName() string {
	return "MarkdownGetByMarkdownID"
}

func (t *MarkdownGetByMarkdownIDEntity) TplType() string {
	return "sql_select"
}

type MarkdownGetAllByMarkdownIDListEntity struct {
	MarkdownIDList []string //

	torm.VolumeMap
}

func (t *MarkdownGetAllByMarkdownIDListEntity) TplName() string {
	return "MarkdownGetAllByMarkdownIDList"
}

func (t *MarkdownGetAllByMarkdownIDListEntity) TplType() string {
	return "sql_select"
}

type MarkdownPaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *MarkdownPaginateWhereEntity) TplName() string {
	return "MarkdownPaginateWhere"
}

func (t *MarkdownPaginateWhereEntity) TplType() string {
	return "text"
}

type MarkdownPaginateTotalEntity struct {
	MarkdownPaginateWhereEntity //

	torm.VolumeMap
}

func (t *MarkdownPaginateTotalEntity) TplName() string {
	return "MarkdownPaginateTotal"
}

func (t *MarkdownPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type MarkdownPaginateEntity struct {
	Limit int //

	MarkdownPaginateWhereEntity //

	Offset int //

	torm.VolumeMap
}

func (t *MarkdownPaginateEntity) TplName() string {
	return "MarkdownPaginate"
}

func (t *MarkdownPaginateEntity) TplType() string {
	return "sql_select"
}

type MarkdownInsertEntity struct {
	APIID string //api 文档标识

	Content string //

	Markdown string //

	MarkdownID string //markdown 文档标识

	Name string //名称(冗余local.en)

	OwnerID int //作者ID

	OwnerName string //作者名称

	ServiceID string //服务标识

	Title string //名称(冗余local.zh)

	torm.VolumeMap
}

func (t *MarkdownInsertEntity) TplName() string {
	return "MarkdownInsert"
}

func (t *MarkdownInsertEntity) TplType() string {
	return "sql_insert"
}

type MarkdownUpdateEntity struct {
	APIID string //api 文档标识

	Content string //

	Markdown string //

	MarkdownID string //markdown 文档标识

	Name string //名称(冗余local.en)

	OwnerID int //作者ID

	OwnerName string //作者名称

	ServiceID string //服务标识

	Title string //名称(冗余local.zh)

	torm.VolumeMap
}

func (t *MarkdownUpdateEntity) TplName() string {
	return "MarkdownUpdate"
}

func (t *MarkdownUpdateEntity) TplType() string {
	return "sql_update"
}

type MarkdownDelEntity struct {
	MarkdownID string //markdown 文档标识

	torm.VolumeMap
}

func (t *MarkdownDelEntity) TplName() string {
	return "MarkdownDel"
}

func (t *MarkdownDelEntity) TplType() string {
	return "sql_update"
}

type ServerGetByServerIDEntity struct {
	ServerID string //服务标识

	torm.VolumeMap
}

func (t *ServerGetByServerIDEntity) TplName() string {
	return "ServerGetByServerID"
}

func (t *ServerGetByServerIDEntity) TplType() string {
	return "sql_select"
}

type ServerGetAllByServerIDListEntity struct {
	ServerIDList []string //

	torm.VolumeMap
}

func (t *ServerGetAllByServerIDListEntity) TplName() string {
	return "ServerGetAllByServerIDList"
}

func (t *ServerGetAllByServerIDListEntity) TplType() string {
	return "sql_select"
}

type ServerPaginateWhereEntity struct {
	torm.VolumeMap
}

func (t *ServerPaginateWhereEntity) TplName() string {
	return "ServerPaginateWhere"
}

func (t *ServerPaginateWhereEntity) TplType() string {
	return "text"
}

type ServerPaginateTotalEntity struct {
	ServerPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ServerPaginateTotalEntity) TplName() string {
	return "ServerPaginateTotal"
}

func (t *ServerPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type ServerPaginateEntity struct {
	Limit int //

	Offset int //

	ServerPaginateWhereEntity //

	torm.VolumeMap
}

func (t *ServerPaginateEntity) TplName() string {
	return "ServerPaginate"
}

func (t *ServerPaginateEntity) TplType() string {
	return "sql_select"
}

type ServerInsertEntity struct {
	Description string //简介

	ExtensionIds string //扩展字段

	Proxy string //代理地址

	ServerID string //服务标识

	ServiceID string //服务标识

	URL string //服务器地址

	torm.VolumeMap
}

func (t *ServerInsertEntity) TplName() string {
	return "ServerInsert"
}

func (t *ServerInsertEntity) TplType() string {
	return "sql_insert"
}

type ServerUpdateEntity struct {
	Description string //简介

	ExtensionIds string //扩展字段

	Proxy string //代理地址

	ServerID string //服务标识

	ServiceID string //服务标识

	URL string //服务器地址

	torm.VolumeMap
}

func (t *ServerUpdateEntity) TplName() string {
	return "ServerUpdate"
}

func (t *ServerUpdateEntity) TplType() string {
	return "sql_update"
}

type ServerDelEntity struct {
	ServerID string //服务标识

	torm.VolumeMap
}

func (t *ServerDelEntity) TplName() string {
	return "ServerDel"
}

func (t *ServerDelEntity) TplType() string {
	return "sql_update"
}
