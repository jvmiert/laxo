import Link from "next/link";
import { useIntl } from "react-intl";

export default function UserMenu() {
  const t = useIntl();

  return (
    <footer
      className="container mx-auto px-4 py-4"
      aria-labelledby="footer-heading"
    >
      <h2 id="footer-heading" className="sr-only">
        Footer
      </h2>
      <div className="mt-12 border-t border-gray-200 pt-2">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-16">
          <div className="xl:grid xl:grid-cols-3 xl:gap-8">
            <div className="text-white xl:col-span-1">
              <p className="tracking-relaxed transform text-lg font-bold tracking-tighter text-indigo-500 transition duration-500 ease-in-out lg:pr-8">
                Laxo
              </p>
              <p className="mt-2 w-1/2 text-sm text-gray-500">
                L17-11, Tầng 17, Tòa nhà Vincom Center, 72 Lê Thánh Tôn, Phường
                Bến Nghé, Quận 1, HCM
              </p>
              <p className="mt-2 w-1/2 text-sm text-gray-500">
                support@laxo.vn
              </p>
            </div>
            <div className="mt-12 grid grid-cols-2 gap-8 xl:col-span-2 xl:mt-0">
              <div className="md:grid md:grid-cols-2 md:gap-8">
                <div>
                  <h3 className="text-sm font-bold uppercase tracking-wider text-indigo-500">
                    {t.formatMessage({
                      defaultMessage: "Navigation",
                      description: "Footer: Navigation header",
                    })}
                  </h3>
                  <ul role="list" className="mt-4 space-y-2">
                    <li>
                      <Link href="/terms">
                        <a className="text-base font-normal text-gray-500 hover:text-indigo-600">
                          {t.formatMessage({
                            defaultMessage: "Terms of Service",
                            description:
                              "Footer: Terms of Service navigation title",
                          })}
                        </a>
                      </Link>
                    </li>
                    <li>
                      <Link href="/privacy">
                        <a className="text-base font-normal text-gray-500 hover:text-indigo-600">
                          {t.formatMessage({
                            defaultMessage: "Privacy Policy",
                            description:
                              "Footer: Privacy Policy navigation title",
                          })}
                        </a>
                      </Link>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
        <p className="text-base text-gray-400 xl:text-center">
          &copy; 2022 Công ty TNHH Laxo.
        </p>
      </div>
    </footer>
  );
}
