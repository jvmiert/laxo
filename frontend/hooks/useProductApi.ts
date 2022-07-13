import { useMemo, useCallback } from "react";
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

export type ChangeSyncRequest = {
  productID: string;
  platform: string;
  state: boolean;
};

export type ChangeImageOrderItem = {
  assetID: string;
  order: number;
};

export default function useProductApi(): {
  doChangeImageOrder: (
    productID: string,
    imageList: Array<ChangeImageOrderItem>,
  ) => Promise<boolean>;
  doGetAssetRank: (assetID: string) => Promise<number>;
  doCreateAsset: (request: CreateAssetRequest) => Promise<CreateAssetReply>;
  doUploadAsset: (assetID: string, file: File) => Promise<UploadAssetReply>;
  doAssetRequest: (request: AssignAssetRequest) => Promise<AssetRequestReply>;
  doChangePlatformSync: (request: ChangeSyncRequest) => Promise<boolean>;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const doChangeImageOrder = useCallback(
    async (
      productID: string,
      imageList: Array<ChangeImageOrderItem>,
    ): Promise<boolean> => {
      try {
        await axiosClient.post(`/change-image-order/${productID}`, {
          assets: imageList,
        });
        return true;
      } catch (error) {
        return false;
      }
    },
    [axiosClient],
  );

  const doChangePlatformSync = useCallback(
    async (r: ChangeSyncRequest): Promise<boolean> => {
      try {
        await axiosClient.post(`/change-platform-sync/${r.productID}`, {
          platform: r.platform,
          state: r.state,
        });
        return true;
      } catch (error) {
        return false;
      }
    },
    [axiosClient],
  );

  const doCreateAsset = useCallback(
    async (request: CreateAssetRequest) => {
      try {
        const res = await axiosClient.post("/asset/create", { ...request });
        return { asset: res.data.asset, upload: res.data.upload, error: false };
      } catch (error) {
        return { asset: undefined, upload: false, error: true };
      }
    },
    [axiosClient],
  );

  const doAssetRequest = useCallback(
    async (request: AssignAssetRequest) => {
      try {
        const res = await axiosClient.post("/asset/manage-product", {
          ...request,
        });
        mutate(`/product/${request.productID}`);
        return { error: false };
      } catch (error) {
        return { error: true };
      }
    },
    [axiosClient, mutate],
  );

  const doUploadAsset = useCallback(
    async (assetID: string, file: File) => {
      try {
        const res = await axiosClient.put(`/asset/${assetID}`, file);
        return { asset: res.data, error: false };
      } catch (error) {
        return { asset: undefined, error: true };
      }
    },
    [axiosClient],
  );

  const doGetAssetRank = useCallback(
    async (assetID: string) => {
      try {
        const res = await axiosClient.get(`/asset/rank/${assetID}`);
        return res.data.rank;
      } catch (error) {
        return 0;
      }
    },
    [axiosClient],
  );

  return useMemo(
    () => ({
      doCreateAsset,
      doUploadAsset,
      doAssetRequest,
      doChangePlatformSync,
      doGetAssetRank,
      doChangeImageOrder,
    }),
    [
      doCreateAsset,
      doUploadAsset,
      doAssetRequest,
      doChangePlatformSync,
      doGetAssetRank,
      doChangeImageOrder,
    ],
  );
}
