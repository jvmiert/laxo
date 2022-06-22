import { useIntl } from "react-intl";
import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";

const ProductDetailsSchema = z.object({
  name: z.string({ required_error: "name_required" }),
  sku: z
    .string()
    .min(4, { message: "too_small_sku" })
    .max(1024, { message: "too_big_sku" }),
  sellingPrice: z
    .number({ invalid_type_error: "sellingPrice_should_number" })
    .positive({ message: "sellingPrice_positive" }),
  costPrice: z.optional(
    z.number({ invalid_type_error: "sellingPrice_should_number" }),
  ),
});

export type ProductDetailsSchemaValues = z.infer<typeof ProductDetailsSchema>;

export default function useProductDetailsApi(): [
  validate: (values: ProductDetailsSchemaValues) => ValidationErrors,
] {
  const t = useIntl();

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

  return [validate];
}
