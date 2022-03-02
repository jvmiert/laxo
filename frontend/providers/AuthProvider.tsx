import { ReactChildren, ReactNode } from "react";
import createSafeContext from "@/lib/useSafeContext";
import { useGetAuth } from "@/hooks/swr-hooks";

export interface AuthConsumerProps {
  auth: boolean;
}

export const [useAuth, Provider] = createSafeContext<AuthConsumerProps>();

export const AuthProvider = ({
  children,
}: {
  children: ReactChildren | ReactNode;
}) => {
  const { auth } = useGetAuth();

  const providerValues: AuthConsumerProps = {
    auth,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
