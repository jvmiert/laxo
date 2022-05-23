import { useRouter } from 'next/router'
import { FormattedMessage } from "react-intl";


function getTitle(platform: string): JSX.Element {
  switch (platform.toLowerCase()) {
    case "lazada":
      return <FormattedMessage
        description="Dashboard specific platform lazada title"
        defaultMessage="Lazada Settings"
      />;
    case "shopee":
      return <FormattedMessage
        description="Dashboard specific platform shopee title"
        defaultMessage="Shopee Settings"
      />;
    case "tiki":
      return <FormattedMessage
        description="Dashboard specific platform tiki title"
        defaultMessage="Tiki Settings"
      />;
    default:
      return <></>;
  }
}

export default function NestedDashboardPlatformTitle() {
  const { query: { platform } } = useRouter();
  return getTitle(platform ? platform.toString() : "")
}
