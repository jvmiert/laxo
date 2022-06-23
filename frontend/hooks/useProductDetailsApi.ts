import axios from "axios";
import { useIntl } from "react-intl";
import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";

import { useAxios } from "@/providers/AxiosProvider";
import { ResponseError, LaxoProductDetailsResponse } from "@/types/ApiResponse";

const ProductDetailsSchema = z.object({
  name: z.string({ required_error: "name_required" }),
  msku: z
    .string()
    .min(4, { message: "too_small_sku" })
    .max(1024, { message: "too_big_sku" }),
  sellingPrice: z
    .number({ invalid_type_error: "sellingPrice_should_number" })
    .positive({ message: "sellingPrice_positive" }),
  costPrice: z.optional(
    z.number({ invalid_type_error: "costPrice_should_number" }),
  ),
});

export type ProductDetailsSchemaValues = z.infer<typeof ProductDetailsSchema>;

type ProductDetailsSubmissionErrors = {
  name?: string;
  msku?: string;
  sellingPrice?: string;
  costPrice?: string;
  [FORM_ERROR]?: string;
};

export default function useProductDetailsApi(
  productID: string,
): [
  validate: (values: ProductDetailsSchemaValues) => ValidationErrors,
  submit: (values: ProductDetailsSchemaValues) => Promise<SubmissionErrors>,
] {
  const t = useIntl();
  const { axiosClient } = useAxios();

  const generalError = {
    [FORM_ERROR]: t.formatMessage({
      defaultMessage:
        "Having trouble saving your product, please try again later",
      description: "Product Details Form: general failure",
    }),
  };

  const submitForm = async (
    values: ProductDetailsSchemaValues,
  ): Promise<SubmissionErrors> => {
    try {
      await axiosClient.post<LaxoProductDetailsResponse>(
        `/product/${productID}`,
        { ...values },
      );
    } catch (error) {
      if (axios.isAxiosError(error) && error.code && error.response) {
        if (parseInt(error.code) >= 500) {
          return generalError;
        }

        const errorObject = error.toJSON() as ResponseError;
        if (Object.keys(errorObject.errorDetails).length == 0) {
          return generalError;
        }

        const errors: ProductDetailsSubmissionErrors = {};
        Object.keys(errorObject.errorDetails).forEach((key) => {
          if (["name", "msku", "sellingPrice", "costPrice"].includes(key)) {
            errors[key as keyof ProductDetailsSubmissionErrors] =
              errorObject.errorDetails[key];
          }
        });
        return errors;
      }

      return generalError;
    }
  };

  const validate = (values: ProductDetailsSchemaValues): ValidationErrors => {
    let errors: { [key: string]: string } = {};

    const validationResult = ProductDetailsSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        let errorMessage: string;
        switch (validation.message) {
          case "name_required":
            errorMessage = t.formatMessage({
              defaultMessage: "Please enter a product name",
              description: "Product Details Form: name validation failure",
            });
            break;
          case "too_small_sku":
            errorMessage = t.formatMessage({
              defaultMessage: "Your SKU should be at least 4 characters",
              description: "Product Details Form: sku min validation failure",
            });
            break;
          case "too_big_sku":
            errorMessage = t.formatMessage({
              defaultMessage:
                "Your SKU should not be more than 1024 characters",
              description: "Product Details Form: sku max validation failure",
            });
            break;
          case "sellingPrice_positive":
            errorMessage = t.formatMessage({
              defaultMessage: "Selling price should be positive",
              description:
                "Product Details Form: price positive validation failure",
            });
            break;
          case "sellingPrice_should_number":
          case "costPrice_should_number":
            errorMessage = t.formatMessage({
              defaultMessage: "Price should be a number",
              description:
                "Product Details Form: price validation validation failure",
            });
            break;
          default:
            errorMessage = t.formatMessage({
              defaultMessage: "Required",
              description: "Product Details Form: general validation failure",
            });
        }
        errors[validation.path[0]] = errorMessage;
      });
    }
    return errors;
  };

  return [validate, submitForm];
}
