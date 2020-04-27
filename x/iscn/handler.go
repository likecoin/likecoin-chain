package iscn

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cbornode "github.com/ipfs/go-ipld-cbor"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateIscn:
			return handleMsgCreateIscn(ctx, msg, keeper)
		case MsgAddEntity:
			return handleMsgAddEntity(ctx, msg, keeper)
		case MsgAddRightTerms:
			return handleMsgAddRightTerms(ctx, msg, keeper)
		default:
			errMsg := fmt.Sprintf("unrecognized iscn message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleEntity(ctx sdk.Context, entity IscnDataField, keeper Keeper) (*CID, error) {
	switch entity.Type() {
	case NestedIscnData:
		e, _ := entity.AsIscnData()
		return keeper.SetEntity(ctx, e)
	case NestedCID:
		cid, _ := entity.AsCID()
		if keeper.GetEntity(ctx, *cid) == nil {
			return nil, fmt.Errorf("unknown entity CID: %s", cid)
		}
		return cid, nil
	default:
		return nil, fmt.Errorf("entity does not match schema")
	}
}

func handleRightTerms(ctx sdk.Context, rightTerms IscnDataField, keeper Keeper) (*CID, error) {
	switch rightTerms.Type() {
	case String:
		rt, _ := rightTerms.AsString()
		return keeper.SetRightTerms(ctx, rt)
	case NestedCID:
		cid, _ := rightTerms.AsCID()
		// Not checking this, allow user to fill in any CIDs which may represent any terms hosted on the IPFS network
		// if keeper.GetRightTerms(ctx, *cid) == nil {
		// 	return nil, fmt.Errorf("unknown right terms CID: %s", cid)
		// }
		return cid, nil
	default:
		return nil, fmt.Errorf("right terms does not match schema")
	}
}

func handleStakeholders(ctx sdk.Context, stakeholders IscnData, keeper Keeper) (*CID, error) {
	stakeholdersArr, _ := stakeholders.Get("stakeholders").AsArray()
	fmt.Printf("handleStakeholders: before anything\n")
	fmt.Printf("handleStakeholders: stakeholders = %v\n", stakeholders)
	for i := 0; i < stakeholdersArr.Len(); i++ {
		stakeholder, _ := stakeholdersArr.Get(i).AsIscnData()
		entityField := stakeholder.Get("stakeholder")
		cid, err := handleEntity(ctx, entityField, keeper)
		if err != nil {
			return nil, err
		}
		stakeholder.Set("stakeholder", *cid)
		n, _ := stakeholder.Get("sharing").AsUint64()
		stakeholder.Set("sharing", uint32(n))
		fmt.Printf("handleStakeholders: After i = %d\n", i)
		fmt.Printf("handleStakeholders: stakeholder = %v\n", stakeholder)
		fmt.Printf("handleStakeholders: stakeholders = %v\n", stakeholders)
	}
	schemaVersion := uint64(1)
	return keeper.SetCidIscnObject(ctx, stakeholders, StakeholdersCodecType, schemaVersion)
}

func handleRights(ctx sdk.Context, rights IscnData, keeper Keeper) (*CID, error) {
	rightsArr, _ := rights.Get("rights").AsArray()
	for i := 0; i < rightsArr.Len(); i++ {
		right, _ := rightsArr.Get(i).AsIscnData()
		holderField := right.Get("holder")
		cid, err := handleEntity(ctx, holderField, keeper)
		if err != nil {
			return nil, err
		}
		right.Set("holder", *cid)

		termsField := right.Get("terms")
		cid, err = handleRightTerms(ctx, termsField, keeper)
		if err != nil {
			return nil, err
		}
		right.Set("terms", *cid)
	}
	schemaVersion := uint64(1)
	return keeper.SetCidIscnObject(ctx, rights, RightsCodecType, schemaVersion)
}

func handleIscnContent(ctx sdk.Context, content IscnDataField, keeper Keeper) (*CID, error) {
	switch content.Type() {
	case NestedIscnData:
		content, _ := content.AsIscnData()
		parent, ok := content.Get("parent").AsCID()
		if ok {
			if keeper.GetEntity(ctx, *parent) == nil {
				return nil, fmt.Errorf("unknown ISCN content parent CID: %s", parent)
			}
			// TODO: parent version check
			content.Set("parent", *parent)
		}
		version, _ := content.Get("version").AsUint64()
		content.Set("version", version)
		return keeper.SetIscnContent(ctx, content)
	case NestedCID:
		cid, _ := content.AsCID()
		if keeper.GetEntity(ctx, *cid) == nil {
			return nil, fmt.Errorf("unknown ISCN content CID: %s", cid)
		}
		return cid, nil
	default:
		return nil, fmt.Errorf("ISCN content does not match schema")
	}
}

func handleMsgCreateIscn(ctx sdk.Context, msg MsgCreateIscn, keeper Keeper) sdk.Result {
	// TODO:
	// 1. store nested fields and construct CIDs
	// 2. validate fields from IscnKernelInput
	// 3. construct IscnKernel
	kernelRawMap := RawIscnMap{}
	err := cbornode.DecodeInto(msg.IscnKernel, &kernelRawMap)
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       fmt.Sprintf("unable to decode ISCN kernel data: %s", err.Error()),
		}
	}
	kernel, ok := KernelSchema.ConstructIscnData(kernelRawMap)
	if !ok {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       "ISCN kernel does not fulfill schema",
		}
	}
	err = keeper.DeductFeeForIscn(ctx, msg.From, msg.IscnKernel)
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	stakeholders, _ := kernel.Get("stakeholders").AsIscnData()
	stakeholdersCID, err := handleStakeholders(ctx, stakeholders, keeper)
	if err != nil {
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	kernel.Set("stakeholders", *stakeholdersCID)
	rights, _ := kernel.Get("rights").AsIscnData()
	rightsCID, err := handleRights(ctx, rights, keeper)
	if err != nil {
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	kernel.Set("rights", *rightsCID)
	content := kernel.Get("content")
	contentCID, err := handleIscnContent(ctx, content, keeper)
	if err != nil {
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	kernel.Set("content", *contentCID)
	// TODO:
	//  - kernel context?
	//  - parent
	//  - check timestamp not later than blocktime
	version, _ := kernel.Get("version").AsUint64()
	kernel.Set("version", version)
	parent := kernel.Get("parent")
	switch parent.Type() {
	case None:
		// New ISCN
		// nil parent case version checking should be handled by ValidateBasic
		_, err = keeper.AddIscnKernel(ctx, kernel)
		if err != nil {
			return sdk.Result{
				Code:      123, // TODO
				Codespace: DefaultCodespace,
				Log:       err.Error(),
			}
		}
	case NestedCID:
		// Old ISCN
		// TODO: check msg.From is the same sender of the old ISCN
		// Need to record ISCN owner
		// TODO: check if the content's parent is pointing to content with the same ISCN ID
		// seems complicated checkings for different weird cases
		parentKernelCID, _ := parent.AsCID()
		parentKernelObj := keeper.GetIscnKernelByCID(ctx, *parentKernelCID)
		if parentKernelObj == nil {
			return sdk.Result{
				Code:      123, // TODO
				Codespace: DefaultCodespace,
				Log:       fmt.Sprintf("unknown parent ISCN kernel CID: %s", parentKernelCID),
			}
		}
		kernel.Set("parent", *parentKernelCID)
		parentVersion, _ := parentKernelObj.GetUint64("version")
		if version != parentVersion+1 {
			return sdk.Result{
				Code:      123, // TODO
				Codespace: DefaultCodespace,
				Log:       "invalid ISCN kernel version",
			}
		}
		iscnID, _ := parentKernelObj.GetBytes("id")
		_, err = keeper.SetIscnKernel(ctx, iscnID, kernel)
		if err != nil {
			return sdk.Result{
				Code:      123, // TODO
				Codespace: DefaultCodespace,
				Log:       err.Error(),
			}
		}
	default:
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       "ISCN kernel parent does not fulfill schema",
		}
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddEntity(ctx sdk.Context, msg MsgAddEntity, keeper Keeper) sdk.Result {
	entityRawMap := RawIscnMap{}
	err := cbornode.DecodeInto(msg.Entity, &entityRawMap)
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       fmt.Sprintf("unable to decode entity data: %s", err.Error()),
		}
	}
	entity, ok := EntitySchema.ConstructIscnData(entityRawMap)
	if !ok {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       "entity does not fulfill schema",
		}
	}
	err = keeper.DeductFeeForIscn(ctx, msg.From, msg.Entity) // TODO: different fee for entity
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	_, err = keeper.SetEntity(ctx, entity)
	if err != nil {
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddRightTerms(ctx sdk.Context, msg MsgAddRightTerms, keeper Keeper) sdk.Result {
	err := keeper.DeductFeeForIscn(ctx, msg.From, []byte(msg.RightTerms)) // TODO: different fee for terms
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      123,
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	_, err = keeper.SetRightTerms(ctx, msg.RightTerms)
	if err != nil {
		return sdk.Result{
			Code:      123, // TODO
			Codespace: DefaultCodespace,
			Log:       err.Error(),
		}
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
