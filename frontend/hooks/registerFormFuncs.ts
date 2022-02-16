import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";
import { useIntl } from "react-intl";

const RegisterSchema = z.object({
  name: z.string({ required_error: "name_required" }),
  email: z
    .string()
    .email({ message: "not_email" })
    .max(300, { message: "too_big_email" }),
  password: z
    .string()
    .regex(/\d/, { message: "contain_digit" })
    .regex(/[^\d]/, { message: "contain_letter" })
    .min(8, { message: "too_small_pw" })
    .max(128, { message: "too_big_pw" }),
});

export type RegisterSchemaValues = z.infer<typeof RegisterSchema>;

export default function useRegisterFuncs(): [
  validate: (values: RegisterSchemaValues) => ValidationErrors,
  submit: (values: RegisterSchemaValues) => Promise<SubmissionErrors>,
] {
  const t = useIntl();

  const submitForm = async (
    values: RegisterSchemaValues,
  ): Promise<SubmissionErrors> => {
    return Promise;
  };

  const validate = (values: RegisterSchemaValues): ValidationErrors => {
    const errors: { [key: string]: string } = {};

    const validationResult = RegisterSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        let errorMessage: string;
        switch (validation.message) {
          case "name_required":
            errorMessage = t.formatMessage({
              defaultMessage: "Please enter your full name",
              description: "Register Form: name validation failure",
            });
            break;

          case "not_email":
            errorMessage = t.formatMessage({
              defaultMessage: "Please enter a valid email address",
              description: "Register Form: email validation failure",
            });
            break;

          case "too_big_email":
            errorMessage = t.formatMessage({
              defaultMessage: "Email should be less than 300 characters",
              description: "Register Form: email max length validation failure",
            });
            break;

          case "too_small_pw":
            errorMessage = t.formatMessage({
              defaultMessage: "Password should be at least 8 characters",
              description:
                "Register Form: password min length validation failure",
            });
            break;

          case "too_big_pw":
            errorMessage = t.formatMessage({
              defaultMessage: "Password should be less than 128 characters",
              description:
                "Register Form: password min length validation failure",
            });
            break;

          case "contain_digit":
            errorMessage = t.formatMessage({
              defaultMessage: "Password must contain a digit",
              description:
                "Register Form: password contain digit validation failure",
            });
            break;

          case "contain_letter":
            errorMessage = t.formatMessage({
              defaultMessage: "Password must contain a letter",
              description:
                "Register Form: password contain letter validation failure",
            });
            break;

          default:
            errorMessage = t.formatMessage({
              defaultMessage: "Required",
              description: "Register Form: general validation failure",
            });
        }
        errors[validation.path[0]] = errorMessage;
      });
    }
    return errors;
  };

  return [validate, submitForm];
}
