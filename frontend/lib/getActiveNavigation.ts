export const getDashboardActiveItem = (
  currentPath: string,
  targetPath: string,
): boolean => {
  let active = false;
  const splitCurrentPath = currentPath.split("/");
  const splitTargetPath = targetPath.split("/");

  const currentFirstLevel =
    splitCurrentPath.length >= 3 ? splitCurrentPath[2] : "";
  const targetFirstLevel =
    splitTargetPath.length >= 3 ? splitTargetPath[2] : "";

  const currentSecondLevel =
    splitCurrentPath.length >= 4 ? splitCurrentPath[3] : "";
  const targetSecondLevel =
    splitTargetPath.length >= 4 ? splitTargetPath[3] : "";

  switch (targetFirstLevel) {
    case "platforms":
      if (targetSecondLevel != "" || currentSecondLevel != "") {
        active = targetSecondLevel == currentSecondLevel;
      } else {
        active = targetFirstLevel == currentFirstLevel;
      }
      break;
    default:
      active = targetFirstLevel == currentFirstLevel;
      break;
  }

  return active;
};
