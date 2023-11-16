import Loader from "@/components/common/loader";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
	doGoogleLogin,
	selectGoogleLoginIsLoading,
} from "@/store/reducers/user";
import {
	ErrorCode,
	GoogleCredentialResponse,
	GoogleLogin,
	GoogleLoginProps,
	GoogleOAuthProvider,
} from "@react-oauth/google";

export default function GoogleLoginContainer() {
	const loading = useAppSelector(selectGoogleLoginIsLoading);
	const dispatch = useAppDispatch();

	const onSignIn = (googleUser: any) => {
		var profile = googleUser.getBasicProfile();
		console.log("ID: " + profile.getId());
		console.log("Name: " + profile.getName());
		console.log("Image URL: " + profile.getImageUrl());
		console.log("Email: " + profile.getEmail());
	};

	const responseMessage = (creds: GoogleCredentialResponse) => {
		const token = creds.credential as string;
		dispatch(doGoogleLogin(token));
	};

	if (!process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID) {
		return;
	}

	return (
		<GoogleOAuthProvider
			clientId={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID}
		>
			<div className="my-8 max-w-xs">
				<h2 className="mt-2 mb-4 text-l text-stone-100">
					With provider:
				</h2>

				{loading ? (
					<div className="absolute top-0 right-0 left-0 bottom-0 width-100 height-100 grid bg-opacity-60 bg-slate-900">
						<div className="place-self-center">
							<Loader />
						</div>
					</div>
				) : (
					<GoogleLogin onSuccess={responseMessage} />
				)}
			</div>
		</GoogleOAuthProvider>
	);
}
