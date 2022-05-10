import { ReactNode } from "react";
import createSafeContext from "@/lib/useSafeContext";
import { useGetAuth } from "@/hooks/swrHooks";

export interface AuthConsumerProps {
  auth: boolean;
}

export const [useAuth, Provider] = createSafeContext<AuthConsumerProps>();

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const { auth } = useGetAuth();

  const providerValues: AuthConsumerProps = {
    auth,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
