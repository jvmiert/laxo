import cc from "classcat";
import { useIntl } from "react-intl";
import { Form, Field } from "react-final-form";
import { SubmissionErrors } from "final-form";
import createDecorator from "final-form-focus";
import type { ReactElement } from "react";
import { useRef, useEffect } from "react";
import { defineMessage } from "react-intl";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";

import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DashboardLayout from "@/components/DashboardLayout";
import useProductDetailsApi, {
  ProductDetailsSchemaValues,
} from "@/hooks/useProductDetailsApi";
import { formatPrice, parsePrice } from "@/lib/priceFormat";
import Editor from "@/components/slate/Editor";
import AssetInsertDialog from "@/components/dashboard/product/AssetInsertDialog";
import AssetManagement from "@/components/dashboard/product/AssetManagement/AssetManagement";
import { useDashboard } from "@/providers/DashboardProvider";

const focusOnError = createDecorator<ProductDetailsSchemaValues>();

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardNewProductProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardNewProduct(props: DashboardNewProductProps) {
  const t = useIntl();
  const [validate, _, submitCreate] = useProductDetailsApi("");

  const { slateEditorRef, productAssetListRef } = useDashboard();
  const initialFocusRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    initialFocusRef.current && initialFocusRef.current.focus();
  }, []);

  const submitFunc = async (
    values: ProductDetailsSchemaValues,
  ): Promise<SubmissionErrors> => {
    if (slateEditorRef.current) {
      values["description"] = slateEditorRef.current.children;
    }

    const assets = productAssetListRef.current
      ? productAssetListRef.current
      : [];

    const result = submitCreate({
      ...values,
      assets,
    });

    return {};
  };

  return (
    <div className="mx-auto max-w-5xl">
      <AssetInsertDialog />
      <div className="space-y-3">
        <div className="rounded-md bg-white py-7 px-6 shadow-sm">
          <Form
            onSubmit={submitFunc}
            validate={validate}
            validateOnBlur
            decorators={[focusOnError]}
            render={({ handleSubmit, submitError, submitting }) => (
              <form
                onSubmit={handleSubmit}
                id="generalEditForm"
                className="grid grid-cols-8 gap-4"
              >
                <div className="col-span-8">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    General
                  </h3>
                </div>

                <div className="col-span-5">
                  <label
                    className="mb-1 block text-sm text-gray-700"
                    htmlFor="name"
                  >
                    {t.formatMessage({
                      description:
                        "General product management: form name label",
                      defaultMessage: "Name",
                    })}
                  </label>
                  <Field<string>
                    autoFocus
                    name="name"
                    render={({ input, meta }) => {
                      const attemped = !meta.pristine || meta.submitFailed;
                      const showError =
                        !meta.dirtySinceLastSubmit &&
                        (!meta.active || meta.submitFailed) &&
                        attemped &&
                        meta.touched &&
                        (meta.error || meta.submitError) &&
                        !meta.submitting;

                      return (
                        <>
                          <input
                            autoFocus
                            ref={initialFocusRef}
                            {...input}
                            className={cc([
                              "focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200",
                              { "border-red-500": showError },
                            ])}
                            type="text"
                            placeholder={t.formatMessage({
                              description:
                                "General product management: form name placeholder",
                              defaultMessage: "Your product name",
                            })}
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
                  <label
                    className="mb-1 block text-sm text-gray-700"
                    htmlFor="msku"
                  >
                    {t.formatMessage({
                      description: "General product management: form sku label",
                      defaultMessage: "SKU",
                    })}
                  </label>
                  <Field<string>
                    name="msku"
                    render={({ input, meta }) => {
                      const attemped = !meta.pristine || meta.submitFailed;
                      const showError =
                        (!meta.active ||
                          meta.submitFailed ||
                          !meta.dirtySinceLastSubmit) &&
                        attemped &&
                        meta.touched &&
                        (meta.error || meta.submitError) &&
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
                            placeholder={t.formatMessage({
                              description:
                                "General product management: form sku placeholder",
                              defaultMessage: "your-unique-product-sku-123",
                            })}
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
                  <label
                    className="mb-1 block text-sm text-gray-700"
                    htmlFor="name"
                  >
                    {t.formatMessage({
                      description:
                        "General product management: form selling price label",
                      defaultMessage: "Selling Price",
                    })}
                  </label>
                  <Field<number, HTMLInputElement, string>
                    name="sellingPrice"
                    format={formatPrice}
                    parse={parsePrice}
                    render={({ input, meta }) => {
                      const attemped = !meta.pristine || meta.submitFailed;
                      const showError =
                        (!meta.active ||
                          meta.submitFailed ||
                          !meta.dirtySinceLastSubmit) &&
                        attemped &&
                        meta.touched &&
                        (meta.error || meta.submitError) &&
                        !meta.submitting;

                      return (
                        <>
                          <div className="flex rounded shadow">
                            <input
                              {...input}
                              className={cc([
                                "focus:shadow-outline w-full appearance-none rounded-tl rounded-bl border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200",
                                {
                                  "border-y-red-500 border-l-red-500":
                                    showError,
                                },
                              ])}
                              type="text"
                              placeholder={t.formatMessage({
                                description:
                                  "General product management: form price placeholder",
                                defaultMessage: "235.000",
                              })}
                            />
                            <span
                              className={cc([
                                "inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500",
                                {
                                  "border-y-red-500 border-r-red-500":
                                    showError,
                                },
                              ])}
                            >
                              ₫
                            </span>
                          </div>
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
                <div className="col-start-6 col-end-9">
                  <label
                    className="mb-1 block text-sm text-gray-700"
                    htmlFor="name"
                  >
                    {t.formatMessage({
                      description:
                        "General product management: form cost price label",
                      defaultMessage: "Cost Price",
                    })}
                  </label>
                  <Field<number, HTMLInputElement, string>
                    name="costPrice"
                    format={formatPrice}
                    parse={parsePrice}
                    render={({ input, meta }) => {
                      const attemped = !meta.pristine || meta.submitFailed;
                      const showError =
                        (!meta.active ||
                          meta.submitFailed ||
                          !meta.dirtySinceLastSubmit) &&
                        attemped &&
                        meta.touched &&
                        (meta.error || meta.submitError) &&
                        !meta.submitting;

                      return (
                        <>
                          <div className="flex rounded shadow">
                            <input
                              {...input}
                              className={cc([
                                "focus:shadow-outline w-full appearance-none rounded-tl rounded-bl border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200",
                                {
                                  "border-y-red-500 border-l-red-500":
                                    showError,
                                },
                              ])}
                              type="text"
                              placeholder={t.formatMessage({
                                description:
                                  "General product management: form cost price placeholder",
                                defaultMessage: "135.000",
                              })}
                            />
                            <span
                              className={cc([
                                "inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500",
                                {
                                  "border-y-red-500 border-r-red-500":
                                    showError,
                                },
                              ])}
                            >
                              ₫
                            </span>
                          </div>
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
                  <Editor />
                </div>
                {submitError && (
                  <p className="mb-2 text-xs italic text-red-500">
                    {submitError}
                  </p>
                )}

                <div className="col-span-8">
                  <h3 className="pt-8 text-lg font-medium leading-6 text-gray-900">
                    {t.formatMessage({
                      defaultMessage: "Media",
                      description: "Product detail edit: media title",
                    })}
                  </h3>
                </div>
                <div className="col-span-8">
                  <AssetManagement mediaList={[]} />
                </div>

                <div className="col-span-8">
                  <button
                    disabled={submitting}
                    className="mt-8 w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
                    type="submit"
                  >
                    {t.formatMessage({
                      defaultMessage: "Create",
                      description: "Create Product Page: Create Button",
                    })}
                  </button>
                </div>
              </form>
            )}
          />
        </div>
      </div>
    </div>
  );
}

DashboardNewProduct.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard new product title",
        defaultMessage: "New Product",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(DashboardNewProduct);
