import { useIntl } from "react-intl";

export default function NotificationEmpty() {
  const t = useIntl();
  return (
    <div className="text-md m-5 flex rounded-lg border-2 border-dashed border-slate-400 p-5 text-center text-slate-400">
      {t.formatMessage({
        defaultMessage: "You don't have any notifications yet",
        description: "Notification: empty",
      })}
    </div>
  );
}
