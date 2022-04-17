import { grpc } from "@improbable-eng/grpc-web";
import axios from "axios";
import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { ResponseError } from "@/types/ApiResponse";
import { ProductService } from "@/proto/product_pb_service";
import { ProductRetrieveUpdateRequest } from "@/proto/product_pb";

const myTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export default function useShopApi(): {
  doCreateShop: (shopName: string) => Promise<ResponseError>;
  getProductRetrieveStatusUpdate: (retrieveID: string) => void;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const getProductRetrieveStatusUpdate = async (retrieveID: string) => {
    const retrieveProduct = new ProductRetrieveUpdateRequest();
    retrieveProduct.setShopid("cool_shop");
    retrieveProduct.setRetrieveid(retrieveID);

    grpc.invoke(ProductService.GetProductRetrieveUpdate, {
      request: retrieveProduct,
      host: "http://localhost:8081",
      transport: myTransport,
      onMessage: (res) => {
        console.log("onMessage", res);
      },
      onEnd: (res) => {
        console.log("onEnd", res);
      },
    });
  };

  const doCreateShop = async (shopName: string) => {
    try {
      await axiosClient.post("/shop", { shopName });
    } catch (error) {
      if (axios.isAxiosError(error)) {
        mutate("/shop");
        if (error.response?.data instanceof Object) {
          return {
            success: false,
            error: true,
            errorDetails: error.response.data.errorDetails,
          };
        }
        return { success: false, error: true, errorDetails: {} };
      }

      throw error;
    }

    mutate("/shop");
    return { success: true, error: false, errorDetails: {} };
  };

  return { doCreateShop, getProductRetrieveStatusUpdate };
}
