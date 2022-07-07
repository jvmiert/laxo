import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";
import { useRouter } from "next/router";
import { useIntl } from "react-intl";
import useShopApi from "@/hooks/useShopApi";

const ShopSchema = z.object({
  shopName: z
    .string()
    .min(6, { message: "too_small_name" })
    .max(300, { message: "too_big_shop_name" }),
});

export type ShopSchemaValues = z.infer<typeof ShopSchema>;

export default function useShopFuncs(): [
  validate: (values: ShopSchemaValues) => ValidationErrors,
  submit: (values: ShopSchemaValues) => Promise<SubmissionErrors>,
] {
  const t = useIntl();
  const { push } = useRouter();
  const { doCreateShop } = useShopApi();
  const submitForm = async (
    values: ShopSchemaValues,
  ): Promise<SubmissionErrors> => {
    const { success, error, errorDetails } = await doCreateShop(
      values.shopName,
    );

    if (error) {
      if (Object.keys(errorDetails).length == 0) {
        return {
          [FORM_ERROR]: t.formatMessage({
            defaultMessage: "Create shop unknown error",
            description: "Shop Create Form: general failure",
          }),
        };
      }
      const errors: { [key: string]: string } = {};
      Object.keys(errorDetails).forEach((key) => {
        switch (errorDetails[key].code) {
          case "validation_length_out_of_range":
            errors[key] = t.formatMessage({
              defaultMessage:
                "Make sure your shop name is between 6 and 300 characters",
              description: "Shop Form: length validation failed",
            });
          default:
            errors[key] = t.formatMessage({
              defaultMessage: "Something went wrong, please try again later",
              description: "Shop Form: unknown failure",
            });
            break;
        }
      });
      return errors;
    }

    if (success) {
      push("/setup-shop/connect");
      return {};
    }
  };

  const validate = (values: ShopSchemaValues): ValidationErrors => {
    const errors: { [key: string]: string } = {};

    const validationResult = ShopSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        let errorMessage: string;
        switch (validation.message) {
          case "too_small_name":
            errorMessage = t.formatMessage({
              defaultMessage: "Shop name should be at least 6 characters",
              description:
                "Shop Create Form: shop name min length validation failure",
            });
            break;

          case "too_big_name":
            errorMessage = t.formatMessage({
              defaultMessage: "Shop name should be less than 300 characters",
              description:
                "Shop Create Form: shop name min length validation failure",
            });
            break;

          default:
            errorMessage = t.formatMessage({
              defaultMessage: "Required",
              description: "Shop Create Form: general validation failure",
            });
        }
        errors[validation.path[0]] = errorMessage;
      });
    }
    return errors;
  };

  return [validate, submitForm];
}
