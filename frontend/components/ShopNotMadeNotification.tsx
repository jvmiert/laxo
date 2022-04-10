import { useIntl } from "react-intl";
import { ArrowNarrowRightIcon } from "@heroicons/react/outline";
import Link from "next/link";

export default function ShopNotMadeNotification() {
  const t = useIntl();

  return (
    <p>
      {t.formatMessage({
        defaultMessage:
          "You didn't create your shop yet, please make your shop.",
        description: "No shop notifcation",
      })}
      <br />
      <Link href="/setup-shop/create" passHref>
        <a className="cursor-pointer font-semibold text-indigo-500">
          {t.formatMessage({
            defaultMessage: "Create your store",
            description: "No shop: create store button",
          })}{" "}
          <ArrowNarrowRightIcon className="inline h-4 w-4" />
        </a>
      </Link>
    </p>
  );
}
