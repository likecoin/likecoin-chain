package types

var (
	EventTypeCreateIscn     = "create_iscn"
	EventTypeAddEntity      = "add_entity"
	EventTypeAddRightTerms  = "add_right_terms"
	EventTypeAddIscnContent = "add_iscn_content"
	EventTypeAddIscnKernel  = "add_iscn_kernel"

	AttributeKeyIscnID         = "iscn_id"
	AttributeKeyIscnKernelCid  = "iscn_kernel_cid"
	AttributeKeyIscnContentCid = "iscn_content_cid"
	AttributeKeyEntityCid      = "entity_cid"
	AttributeKeyRightTermsCid  = "right_terms_cid"
	AttributeValueCategory     = ModuleName
)
