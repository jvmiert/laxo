import { useIntl } from "react-intl";

export default function UserMenu() {
  const t = useIntl();

  return (
    <footer className="container mx-auto px-4 py-4">
      <div className="mt-12 border-t border-gray-200 pt-8">
        <p className="text-base text-gray-400 xl:text-center">
          &copy; 2022 Công ty TNHH Laxo.
        </p>
      </div>
    </footer>
  );
}
