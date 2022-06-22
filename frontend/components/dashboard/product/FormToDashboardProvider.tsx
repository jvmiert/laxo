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
  const { values, valid, submitting } = useFormState();
  const { reset } = useForm();

  const {
    productDetailFormResetRef,
    productDetailFormIsDirty,
    toggleProductDetailFormDirtyState,
    productDetailSubmitIsDisabled,
    toggleProductDetailSubmitIsDisabled,
  } = useDashboard();

  const changed = useMemo(
    () => diff(values, initialValues, { cyclesFix: false }).length > 0,
    [initialValues, values],
  );

  useEffect(() => {
    const disabled = !valid || submitting;

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
