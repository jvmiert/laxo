import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";
import { useRouter } from "next/router";
import { useIntl } from "react-intl";
import useLoginApi from "@/hooks/useLoginApi";
import useRedirectSafely from "@/hooks/redirectSafely";

const LoginSchema = z.object({
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

export type LoginSchemaValues = z.infer<typeof LoginSchema>;

export default function useLoginFuncs(): [
  validate: (values: LoginSchemaValues) => ValidationErrors,
  submit: (values: LoginSchemaValues) => Promise<SubmissionErrors>,
] {
  const t = useIntl();
  const { doLogin } = useLoginApi();
  const { query } = useRouter();
  const { redirectSafely } = useRedirectSafely();
  const submitForm = async (
    values: LoginSchemaValues,
  ): Promise<SubmissionErrors> => {
    const { success, error, errorDetails } = await doLogin(
      values.email,
      values.password,
    );

    if (error) {
      if (Object.keys(errorDetails).length == 0) {
        return {
          [FORM_ERROR]: t.formatMessage({
            defaultMessage: "Login Failed",
            description: "Login Form: general failure",
          }),
        };
      }
      const errors: { [key: string]: string } = {};
      Object.keys(errorDetails).forEach((key) => {
        errors[key] = errorDetails[key];
      });
      return errors;
    }

    if (success) {
      redirectSafely(query?.next ? (query.next as string) : "/");
      return {};
    }
  };

  const validate = (values: LoginSchemaValues): ValidationErrors => {
    const errors: { [key: string]: string } = {};

    const validationResult = LoginSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        let errorMessage: string;
        switch (validation.message) {
          case "not_email":
            errorMessage = t.formatMessage({
              defaultMessage: "Please enter a valid email address",
              description: "Login Form: email validation failure",
            });
            break;

          case "too_big_email":
            errorMessage = t.formatMessage({
              defaultMessage: "Email should be less than 300 characters",
              description: "Login Form: email max length validation failure",
            });
            break;

          case "too_small_pw":
            errorMessage = t.formatMessage({
              defaultMessage: "Password should be at least 8 characters",
              description: "Login Form: password min length validation failure",
            });
            break;

          case "too_big_pw":
            errorMessage = t.formatMessage({
              defaultMessage: "Password should be less than 128 characters",
              description: "Login Form: password min length validation failure",
            });
            break;

          case "contain_digit":
            errorMessage = t.formatMessage({
              defaultMessage: "Password must contain a digit",
              description:
                "Login Form: password contain digit validation failure",
            });
            break;

          case "contain_letter":
            errorMessage = t.formatMessage({
              defaultMessage: "Password must contain a letter",
              description:
                "Login Form: password contain letter validation failure",
            });
            break;

          default:
            errorMessage = t.formatMessage({
              defaultMessage: "Required",
              description: "Login Form: general validation failure",
            });
        }
        errors[validation.path[0]] = errorMessage;
      });
    }
    return errors;
  };

  return [validate, submitForm];
}
