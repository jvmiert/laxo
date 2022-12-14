import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";
import { useRouter } from "next/router";
import { useIntl } from "react-intl";
import useRegisterApi from "@/hooks/registerUser";
import useRedirectSafely from "@/hooks/redirectSafely";

const RegisterSchema = z.object({
  fullname: z.string({ required_error: "name_required" }),
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
  const [doRegistration] = useRegisterApi();
  const { query } = useRouter();
  const { redirectSafely } = useRedirectSafely();

  const submitForm = async (
    values: RegisterSchemaValues,
  ): Promise<SubmissionErrors> => {
    const { success, error, errorDetails } = await doRegistration(
      values.email,
      values.password,
      values.fullname,
    );

    if (error) {
      if (Object.keys(errorDetails).length == 0) {
        return {
          [FORM_ERROR]: t.formatMessage({
            defaultMessage: "Registration Failed",
            description: "Registration Form: general failure",
          }),
        };
      }
      const errors: { [key: string]: string } = {};
      Object.keys(errorDetails).forEach((key) => {
        switch (errorDetails[key].code) {
          //@TODO: Add the additional validations that are already performed by zod
          //       but are double checked in the backend with Ozzo
          case "already_exists":
            errors[key] = t.formatMessage({
              defaultMessage: "User already exists",
              description: "Register Form: user already exists",
            });
            break;
          default:
            errors[key] = t.formatMessage({
              defaultMessage: "Something went wrong, please try again later",
              description: "Register Form: unknown failure",
            });
            break;
        }
      });
      return errors;
    }

    if (success) {
      redirectSafely(query?.next ? (query.next as string) : "/dashboard/home");
      return {};
    }
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
