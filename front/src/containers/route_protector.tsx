import { useAppSelector } from "@/store/hooks";
import { selectIsLoggedIn } from "@/store/reducers/user";
import { useRouter } from "next/router";

export default function RouteProtector({ children }: any) {
	const router = useRouter();
	const authenticated = useAppSelector(selectIsLoggedIn);

	const nonAuthRoutes = ["/", "/login", "/register"];
	const isNonAuthRoute = nonAuthRoutes.find(
		(route) => route === window.location.pathname
	);

	const isAuthRoute = window.location.pathname.startsWith("/polls");

	if (authenticated && isNonAuthRoute) {
		router.push("/polls");
	} else if (!authenticated && isAuthRoute) {
		router.push("/");
	}

	return children;
}
