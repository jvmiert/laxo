import cc from "classcat";
import { Form, Field } from "react-final-form";
import { SubmissionErrors } from "final-form";
import createDecorator from "final-form-focus";
import { useIntl } from "react-intl";

import Editor from "@/components/slate/Editor";
import FormToDashboardProvider from "@/components/dashboard/product/FormToDashboardProvider";
import { LaxoProduct } from "@/types/ApiResponse";
import useProductDetailsApi, {
  ProductDetailsSchemaValues,
} from "@/hooks/useProductDetailsApi";
import { useGetLaxoProductDetails } from "@/hooks/swrHooks";
import { useDashboard } from "@/providers/DashboardProvider";
import { formatPrice, parsePrice } from "@/lib/priceFormat";

const focusOnError = createDecorator<ProductDetailsSchemaValues>();

export type GeneralEditProps = {
  product: LaxoProduct["product"];
};

export default function DetailsGeneralEdit({ product }: GeneralEditProps) {
  const initialValues = {
    name: product.name,
    sellingPrice: parseFloat(
      `${product.sellingPrice.Int}e${product.sellingPrice.Exp}`,
    ),
    costPrice:
      parseFloat(`${product.costPrice.Int}e${product.costPrice.Exp}`) || 0,
    msku: product.msku,
  };

  const { mutate } = useGetLaxoProductDetails(product.id);
  const {
    slateEditorRef,
    dashboardDispatch,
    slateIsDirty,
    toggleSlateDirtyState,
  } = useDashboard();
  const t = useIntl();
  const [validate, submit] = useProductDetailsApi(product.id);

  const submitFunc = async (
    values: ProductDetailsSchemaValues,
  ): Promise<SubmissionErrors> => {
    // The slate editor is not managed by final-form so we add the values now
    if (slateEditorRef.current) {
      values["description"] = slateEditorRef.current.children;
    }
    const result = await submit(values);

    //@TODO: Handle this
    if (!result) return {};

    const [errors, newProduct] = result;
    // Leaving the description value in here causes the final form inital values
    // to always be different due to not managing description with final form
    delete values["description"];

    if (!errors) return {};

    if (Object.keys(errors).length == 0) {
      //@TODO: Figure out how to use the newly returned product to
      //       optimistically update with mutate.
      mutate();

      dashboardDispatch({
        type: "alert",
        alert: {
          type: "success",
          message: t.formatMessage({
            description: "General product management successful save",
            defaultMessage: "Successfully updated your product",
          }),
        },
      });

      if (slateIsDirty) {
        toggleSlateDirtyState();
      }
    }

    return errors;
  };

  return (
    <Form
      onSubmit={submitFunc}
      validate={validate}
      decorators={[focusOnError]}
      initialValues={initialValues}
      render={({ handleSubmit, submitError }) => (
        <form
          onSubmit={handleSubmit}
          id="generalEditForm"
          className="grid grid-cols-8 gap-4"
        >
          <FormToDashboardProvider initialValues={initialValues} />
          <div className="col-span-5">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              {t.formatMessage({
                description: "General product management: form name label",
                defaultMessage: "Name",
              })}
            </label>
            <Field<string>
              name="name"
              render={({ input, meta }) => {
                const attemped = !meta.pristine || meta.submitFailed;
                const unchangedAfterSubmit =
                  meta.submitError && !meta.dirtySinceLastSubmit;
                const showError =
                  attemped &&
                  meta.touched &&
                  (meta.error || unchangedAfterSubmit) &&
                  !meta.submitting;

                return (
                  <>
                    <input
                      {...input}
                      className={cc([
                        "focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200",
                        { "border-red-500": showError },
                      ])}
                      type="text"
                    />
                    {showError && (
                      <span className="text-xs italic text-red-500">
                        {meta.error || meta.submitError}
                      </span>
                    )}
                  </>
                );
              }}
            />
          </div>
          <div className="col-span-3">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="msku">
              {t.formatMessage({
                description: "General product management: form sku label",
                defaultMessage: "SKU",
              })}
            </label>
            <Field<string>
              name="msku"
              render={({ input, meta }) => {
                const attemped = !meta.pristine || meta.submitFailed;
                const unchangedAfterSubmit =
                  meta.submitError && !meta.dirtySinceLastSubmit;
                const showError =
                  attemped &&
                  meta.touched &&
                  (meta.error || unchangedAfterSubmit) &&
                  !meta.submitting;

                return (
                  <>
                    <input
                      {...input}
                      className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                      type="text"
                    />
                    {showError && (
                      <span className="text-xs italic text-red-500">
                        {meta.error || meta.submitError}
                      </span>
                    )}
                  </>
                );
              }}
            />
          </div>
          <div className="col-start-1 col-end-4">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              {t.formatMessage({
                description:
                  "General product management: form selling price label",
                defaultMessage: "Selling Price",
              })}
            </label>
            <div className="flex rounded shadow">
              <Field<number, HTMLInputElement, string>
                name="sellingPrice"
                format={formatPrice}
                parse={parsePrice}
                render={({ input, meta }) => {
                  const attemped = !meta.pristine || meta.submitFailed;
                  const unchangedAfterSubmit =
                    meta.submitError && !meta.dirtySinceLastSubmit;
                  const showError =
                    attemped &&
                    meta.touched &&
                    (meta.error || unchangedAfterSubmit) &&
                    !meta.submitting;

                  return (
                    <>
                      <input
                        {...input}
                        className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                        type="text"
                      />
                      {showError && (
                        <span className="text-xs italic text-red-500">
                          {meta.error || meta.submitError}
                        </span>
                      )}
                    </>
                  );
                }}
              />
              <span className="inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500">
                ₫
              </span>
            </div>
          </div>
          <div className="col-start-6 col-end-9">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              {t.formatMessage({
                description:
                  "General product management: form cost price label",
                defaultMessage: "Cost Price",
              })}
            </label>
            <div className="flex rounded shadow">
              <Field<number, HTMLInputElement, string>
                name="costPrice"
                format={formatPrice}
                parse={parsePrice}
                render={({ input, meta }) => {
                  const attemped = !meta.pristine || meta.submitFailed;
                  const unchangedAfterSubmit =
                    meta.submitError && !meta.dirtySinceLastSubmit;
                  const showError =
                    attemped &&
                    meta.touched &&
                    (meta.error || unchangedAfterSubmit) &&
                    !meta.submitting;

                  return (
                    <>
                      <input
                        {...input}
                        className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                        type="text"
                      />
                      {showError && (
                        <span className="text-xs italic text-red-500">
                          {meta.error || meta.submitError}
                        </span>
                      )}
                    </>
                  );
                }}
              />
              <span className="inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500">
                ₫
              </span>
            </div>
          </div>
          <div className="col-span-8">
            <label
              className="mb-1 block text-sm text-gray-700"
              htmlFor="description"
            >
              {t.formatMessage({
                description:
                  "General product management: form description label",
                defaultMessage: "Description",
              })}
            </label>
            <Editor initialSchema={product.descriptionSlate} />
          </div>
          {submitError && (
            <p className="mb-2 text-xs italic text-red-500">{submitError}</p>
          )}
        </form>
      )}
    />
  );
}
