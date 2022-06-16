import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";

export type Asset = {
  id: string;
  shopID: string;
  murmurHash: string;
  originalFilename: string;
  extension: string;
  fileSize: number;
  width: number;
  height: number;
};

export type UploadAssetReply = {
  asset: Asset | undefined;
  error: boolean;
};

export type CreateAssetReply = {
  asset: Asset | undefined;
  upload: boolean;
  error: boolean;
};

export type AssetRequestReply = {
  error: boolean;
};

export type CreateAssetRequest = {
  originalName: string;
  size: number;
  width: number;
  height: number;
  hash: string;
};

export type AssignAssetRequest = {
  action: "delete" | "inactive" | "active";
  productID: string;
  assetID: string;
  order: Number;
};

export default function useProductApi(): {
  doCreateAsset: (request: CreateAssetRequest) => Promise<CreateAssetReply>;
  doUploadAsset: (assetID: string, file: File) => Promise<UploadAssetReply>;
  doAssetRequest: (request: AssignAssetRequest) => Promise<AssetRequestReply>;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const doCreateAsset = async (request: CreateAssetRequest) => {
    try {
      const res = await axiosClient.post("/asset/create", { ...request });
      return { asset: res.data.asset, upload: res.data.upload, error: false };
    } catch (error) {
      return { asset: undefined, upload: false, error: true };
    }
  };

  const doAssetRequest = async (request: AssignAssetRequest) => {
    try {
      const res = await axiosClient.post("/asset/manage-product", {
        ...request,
      });
      mutate(`/product/${request.productID}`);
      return { error: false };
    } catch (error) {
      return { error: true };
    }
  };

  const doUploadAsset = async (assetID: string, file: File) => {
    try {
      const res = await axiosClient.put(`/asset/${assetID}`, file);
      return { asset: res.data, error: false };
    } catch (error) {
      return { asset: undefined, error: true };
    }
  };

  return { doCreateAsset, doUploadAsset, doAssetRequest };
}