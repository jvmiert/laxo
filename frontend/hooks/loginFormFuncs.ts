import { z } from "zod";
import { SubmissionErrors, ValidationErrors, FORM_ERROR } from "final-form";
import type { NextRouter } from "next/router";
import useLoginApi from "./loginUser";

const LoginSchema = z.object({
  email: z.string().email().max(300),
  password: z.string().min(8).max(128),
});

type LoginSchemaValues = z.infer<typeof LoginSchema>;

export default function useLoginFuncs(
  router: NextRouter,
): [
  validate: (values: LoginSchemaValues) => ValidationErrors,
  submit: (values: LoginSchemaValues) => Promise<SubmissionErrors>,
] {
  const [doLogin] = useLoginApi();

  const submitForm = async (
    values: LoginSchemaValues,
  ): Promise<SubmissionErrors> => {
    const { success, error, errorDetails } = await doLogin(
      values.email,
      values.password,
    );

    if (error) {
      if (Object.keys(errorDetails).length == 0) {
        return { [FORM_ERROR]: "Login Failed" };
      }
      const errors: { [key: string]: string } = {};
      Object.keys(errorDetails).forEach((key) => {
        errors[key] = errorDetails[key];
      });
      return errors;
    }

    if (success) {
      router.push("/");
      return {};
    }
  };

  const validate = (values: LoginSchemaValues): ValidationErrors => {
    const errors: { [key: string]: string } = {};

    const validationResult = LoginSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        errors[validation.path[0]] = validation.message;
      });
    }
    return errors;
  };

  return [validate, submitForm];
}
