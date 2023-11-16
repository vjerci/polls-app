export default function GoogleLoginDisabled() {
	return (
		<div className="my-4">
			<div className="text-gray-500">
				Google login is at the moment disabled for this app
			</div>

			<div className="text-gray-500">
				To enable google login please edit your .zshrc and insert your
				app key, after it rebuild the app
			</div>
		</div>
	);
}
