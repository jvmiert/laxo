import { useDashboard } from "@/providers/DashboardProvider";
import Alert from "@/components/dashboard/Alert";

export default function AlertContainer() {
  const { dashboardState } = useDashboard();
  return (
    <div className="fixed bottom-4 left-1/2 z-50 flex -translate-x-1/2 flex-col justify-center gap-4">
      {dashboardState.alerts.map((a) => (
        <Alert key={a.id} id={a.id} type={a.type} message={a.message} />
      ))}
    </div>
  );
}
