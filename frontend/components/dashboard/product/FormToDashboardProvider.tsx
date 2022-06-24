import { useMemo, useEffect } from "react";
import diff from "microdiff";
import { useFormState, useForm } from "react-final-form";

import { useDashboard } from "@/providers/DashboardProvider";

export type FormToDashboardProviderProps = {
  initialValues: object;
};

export default function FormToDashboardProvider({
  initialValues,
}: FormToDashboardProviderProps) {
  const { values, valid, submitting, submitFailed } = useFormState();
  const { reset } = useForm();

  const {
    productDetailFormResetRef,
    productDetailFormIsDirty,
    toggleProductDetailFormDirtyState,
    productDetailSubmitIsDisabled,
    toggleProductDetailSubmitIsDisabled,
    productDetailFormIsSubmitting,
    toggleProductDetailFormIsSubmitting,
  } = useDashboard();

  const changed = useMemo(
    () => diff(values, initialValues, { cyclesFix: false }).length > 0,
    [initialValues, values],
  );

  useEffect(() => {
    if (submitting && !productDetailFormIsSubmitting) {
      toggleProductDetailFormIsSubmitting();
    }

    if (!submitting && productDetailFormIsSubmitting) {
      toggleProductDetailFormIsSubmitting();
    }
  }, [
    submitting,
    productDetailFormIsSubmitting,
    toggleProductDetailFormIsSubmitting,
  ]);

  useEffect(() => {
    const disabled = (!valid || submitting) && !submitFailed;

    if (disabled && !productDetailSubmitIsDisabled) {
      toggleProductDetailSubmitIsDisabled();
    }

    if (!disabled && productDetailSubmitIsDisabled) {
      toggleProductDetailSubmitIsDisabled();
    }
  }, [
    valid,
    submitting,
    productDetailSubmitIsDisabled,
    toggleProductDetailSubmitIsDisabled,
    submitFailed,
  ]);

  useEffect(() => {
    if (changed && !productDetailFormIsDirty) {
      toggleProductDetailFormDirtyState();
    }
    if (!changed && productDetailFormIsDirty) {
      toggleProductDetailFormDirtyState();
    }
  }, [changed, productDetailFormIsDirty, toggleProductDetailFormDirtyState]);

  useEffect(() => {
    productDetailFormResetRef.current = reset;
  }, [reset, productDetailFormResetRef]);

  return <></>;
}
