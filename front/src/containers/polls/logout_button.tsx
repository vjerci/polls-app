import Link from "@/components/common/link";
import { logout } from "@/store/reducers/user";
import { useDispatch } from "react-redux";

export default function LogoutButton() {
	const dispatch = useDispatch();

	const onLogoutClick = () => {
		dispatch(logout());
	};

	return (
		<Link className="mr-4 block justify-self-end" onClick={onLogoutClick}>
			Logout
		</Link>
	);
}
