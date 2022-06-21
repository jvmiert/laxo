import { Form, Field } from "react-final-form";

import Editor from "@/components/slate/Editor";
import DetailsChangedNotification from "@/components/dashboard/product/DetailsChangedNotification";
import { LaxoProduct } from "@/types/ApiResponse";

const formatPrice = (value: number, name: string): string => {
  return value.toLocaleString("vi-VN");
};

const parsePrice = (value: string, name: string): number => {
  return parseFloat(value.replaceAll(".", "")) || 0;
};

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
    sku: product.msku,
  };

  //@TODO: type this with Zod
  const submitFunc = (values: object) => {
    console.log(values);
  };

  return (
    <Form
      onSubmit={submitFunc}
      initialValues={initialValues}
      render={({ handleSubmit, submitting, submitError }) => (
        <form
          onSubmit={handleSubmit}
          id="generalEditForm"
          className="grid grid-cols-8 gap-4"
        >
          <DetailsChangedNotification
            initialValues={initialValues}
            submitting={submitting}
          />
          <div className="col-span-5">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              Name
            </label>
            <Field<string>
              name="name"
              render={({ input, meta }) => (
                <input
                  {...input}
                  className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                  type="text"
                />
              )}
            />
          </div>
          <div className="col-span-3">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="msku">
              SKU
            </label>
            <Field<string>
              name="sku"
              render={({ input, meta }) => (
                <input
                  {...input}
                  className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                  type="text"
                />
              )}
            />
          </div>
          <div className="col-start-1 col-end-4">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              Selling Price
            </label>
            <div className="flex rounded shadow">
              <Field<number, HTMLInputElement, string>
                name="sellingPrice"
                format={formatPrice}
                parse={parsePrice}
                render={({ input, meta }) => (
                  <input
                    {...input}
                    className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                    type="text"
                  />
                )}
              />
              <span className="inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500">
                ₫
              </span>
            </div>
          </div>
          <div className="col-start-6 col-end-9">
            <label className="mb-1 block text-sm text-gray-700" htmlFor="name">
              Cost Price
            </label>
            <div className="flex rounded shadow">
              <Field<number, HTMLInputElement, string>
                name="costPrice"
                format={formatPrice}
                parse={parsePrice}
                render={({ input, meta }) => (
                  <input
                    {...input}
                    className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                    type="text"
                  />
                )}
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
              Description
            </label>
            <textarea
              name="description"
              rows={8}
              defaultValue={product.description}
              className="focus:shadow-outline block w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
            />
          </div>
          <div className="col-span-8">
            <Editor initialSchema={product.descriptionSlate} />
          </div>
          <button className="invisible" type="submit" />
        </form>
      )}
    />
  );
}
