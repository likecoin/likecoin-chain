package e2e_test

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func parseEventCreateClass(res sdk.TxResponse) types.EventNewClass {
	actualEvent := types.EventNewClass{}

ParseEventCreateClass:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventNewClass" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "parent_iscn_id_prefix" {
						actualEvent.ParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "parent_account" {
						actualEvent.ParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateClass
			}
		}
	}

	return actualEvent
}

func parseEventUpdateClass(res sdk.TxResponse) types.EventUpdateClass {
	actualEvent := types.EventUpdateClass{}

ParseEventUpdateClass:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventUpdateClass" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "parent_iscn_id_prefix" {
						actualEvent.ParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "parent_account" {
						actualEvent.ParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateClass
			}
		}
	}

	return actualEvent
}

func parseEventMintNFT(res sdk.TxResponse) types.EventMintNFT {
	actualEvent := types.EventMintNFT{}

ParseEventMintNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventMintNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_iscn_id_prefix" {
						actualEvent.ClassParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_account" {
						actualEvent.ClassParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventMintNFT
			}
		}
	}

	return actualEvent
}

func parseEventBurnNFT(res sdk.TxResponse) types.EventBurnNFT {
	actualEvent := types.EventBurnNFT{}

ParseEventBurnNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventBurnNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_iscn_id_prefix" {
						actualEvent.ClassParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_account" {
						actualEvent.ClassParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventBurnNFT
			}
		}
	}

	return actualEvent
}

func parseEventCreateBlindBoxContent(res sdk.TxResponse) types.EventCreateBlindBoxContent {
	actualEvent := types.EventCreateBlindBoxContent{}

ParseEventCreateBlindBoxContent:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventCreateBlindBoxContent" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "content_id" {
						actualEvent.ContentId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_iscn_id_prefix" {
						actualEvent.ClassParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_account" {
						actualEvent.ClassParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateBlindBoxContent
			}
		}
	}

	return actualEvent
}

func parseEventUpdateBlindBoxContent(res sdk.TxResponse) types.EventUpdateBlindBoxContent {
	actualEvent := types.EventUpdateBlindBoxContent{}

ParseEventUpdateBlindBoxContent:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventUpdateBlindBoxContent" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "content_id" {
						actualEvent.ContentId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_iscn_id_prefix" {
						actualEvent.ClassParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_account" {
						actualEvent.ClassParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateBlindBoxContent
			}
		}
	}

	return actualEvent
}

func parseEventDeleteBlindBoxContent(res sdk.TxResponse) types.EventDeleteBlindBoxContent {
	actualEvent := types.EventDeleteBlindBoxContent{}

ParseEventDeleteBlindBoxContent:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventDeleteBlindBoxContent" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "content_id" {
						actualEvent.ContentId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_iscn_id_prefix" {
						actualEvent.ClassParentIscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "class_parent_account" {
						actualEvent.ClassParentAccount = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventDeleteBlindBoxContent
			}
		}
	}

	return actualEvent
}

func parseEventCreateListing(res sdk.TxResponse) types.EventCreateListing {
	actualEvent := types.EventCreateListing{}

ParseEventCreateListing:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventCreateListing" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "seller" {
						actualEvent.Seller = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateListing
			}
		}
	}

	return actualEvent
}

func parseEventUpdateListing(res sdk.TxResponse) types.EventUpdateListing {
	actualEvent := types.EventUpdateListing{}

ParseEventUpdateListing:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventUpdateListing" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "seller" {
						actualEvent.Seller = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateListing
			}
		}
	}
	return actualEvent
}

func parseEventBuyNFT(res sdk.TxResponse) types.EventBuyNFT {
	actualEvent := types.EventBuyNFT{}

ParseEventBuyNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventBuyNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "seller" {
						actualEvent.Seller = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "buyer" {
						actualEvent.Buyer = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "price" {
						price, err := strconv.ParseUint(strings.Trim(attr.Value, "\""), 10, 64)
						if err != nil {
							panic(err)
						}
						actualEvent.Price = price
					}
				}
				break ParseEventBuyNFT
			}
		}
	}
	return actualEvent
}

func parseEventCreateOffer(res sdk.TxResponse) types.EventCreateOffer {
	actualEvent := types.EventCreateOffer{}

ParseEventCreateOffer:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventCreateOffer" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "buyer" {
						actualEvent.Buyer = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateOffer
			}
		}
	}

	return actualEvent
}

func parseEventUpdateOffer(res sdk.TxResponse) types.EventUpdateOffer {
	actualEvent := types.EventUpdateOffer{}

ParseEventUpdateOffer:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventUpdateOffer" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "buyer" {
						actualEvent.Buyer = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateOffer
			}
		}
	}

	return actualEvent
}

func parseEventSellNFT(res sdk.TxResponse) types.EventSellNFT {
	actualEvent := types.EventSellNFT{}

ParseEventSellNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventSellNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nft_id" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "seller" {
						actualEvent.Seller = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "buyer" {
						actualEvent.Buyer = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "price" {
						price, err := strconv.ParseUint(strings.Trim(attr.Value, "\""), 10, 64)
						if err != nil {
							panic(err)
						}
						actualEvent.Price = price
					}
					if attr.Key == "full_pay_to_royalty" {
						fullPayToRoyalty, err := strconv.ParseBool(strings.Trim(attr.Value, "\""))
						if err != nil {
							panic(err)
						}
						actualEvent.FullPayToRoyalty = fullPayToRoyalty
					}
				}
				break ParseEventSellNFT
			}
		}
	}
	return actualEvent
}

func parseEventCreateRoyaltyConfig(res sdk.TxResponse) types.EventCreateRoyaltyConfig {
	actualEvent := types.EventCreateRoyaltyConfig{}

ParseEventCreateRoyaltyConfig:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventCreateRoyaltyConfig" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateRoyaltyConfig
			}
		}
	}
	return actualEvent
}

func parseEventUpdateRoyaltyConfig(res sdk.TxResponse) types.EventUpdateRoyaltyConfig {
	actualEvent := types.EventUpdateRoyaltyConfig{}

ParseEventUpdateRoyaltyConfig:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventUpdateRoyaltyConfig" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateRoyaltyConfig
			}
		}
	}
	return actualEvent
}

func parseEventDeleteRoyaltyConfig(res sdk.TxResponse) types.EventDeleteRoyaltyConfig {
	actualEvent := types.EventDeleteRoyaltyConfig{}

ParseEventDeleteRoyaltyConfig:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likechain.likenft.v1.EventDeleteRoyaltyConfig" {
				for _, attr := range event.Attributes {
					if attr.Key == "class_id" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventDeleteRoyaltyConfig
			}
		}
	}
	return actualEvent
}
