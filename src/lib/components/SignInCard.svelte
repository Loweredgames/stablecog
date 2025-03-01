<script lang="ts">
	import { page } from '$app/stores';
	import Button from '$components/buttons/Button.svelte';
	import NoBgButton from '$components/buttons/NoBgButton.svelte';
	import ErrorLine from '$components/ErrorLine.svelte';
	import IconBack from '$components/icons/IconBack.svelte';
	import IconEmail from '$components/icons/IconEmail.svelte';
	import IconGoogle from '$components/icons/IconGoogle.svelte';
	import Input from '$components/Input.svelte';
	import LL from '$i18n/i18n-svelte';
	import { expandCollapse } from '$ts/animation/transitions';
	import { supabase } from '$ts/constants/supabase';
	import { isTouchscreen } from '$ts/stores/isTouchscreen';
	import type { Provider } from '@supabase/supabase-js';
	import { quadOut } from 'svelte/easing';

	export let redirectTo: string | null = null;
	export let isModal = false;

	let email: string;
	let signInStatus: 'idle' | 'loading' | 'error' | 'sent-otp' = 'idle';
	let errorText: string | null = null;

	async function signIn() {
		if (!email.includes('@')) {
			errorText = $LL.Error.InvalidEmail();
			return;
		}
		signInStatus = 'loading';
		const { data: sData, error: sError } = await supabase.auth.signInWithOtp({
			email,
			options: {
				emailRedirectTo: `${$page.url.origin}/api/auth/callback?redirect_to=${
					redirectTo ? encodeURIComponent(redirectTo) : ''
				}`
			}
		});
		if (sError) {
			console.log(sError);
			signInStatus = 'error';
			if (
				sError.message === 'For security purposes, you can only request this once every 60 seconds'
			) {
				errorText = $LL.Error.OnceEvery60Seconds();
			} else {
				errorText = $LL.Error.SomethingWentWrong();
			}
			return;
		}
		console.log(sData);
		signInStatus = 'sent-otp';
	}

	async function signInWithOAuth(provider: Provider) {
		signInStatus = 'loading';
		const { data: sData, error: sError } = await supabase.auth.signInWithOAuth({
			provider
		});
		if (sError) {
			console.log(sError);
			signInStatus = 'error';
			errorText = $LL.Error.SomethingWentWrong();
			return;
		}
		console.log(sData);
	}
</script>

<div
	class="max-w-full flex flex-col items-center justify-center {isModal
		? 'bg-c-bg ring-2 ring-c-bg-secondary p-5 md:px-10 md:py-7 rounded-3xl shadow-2xl shadow-c-shadow/[var(--o-shadow-strong)]'
		: ''}"
>
	{#if signInStatus === 'sent-otp'}
		<div class="mb-2">
			<IconEmail class="w-20 h-20 text-c-on-bg" />
		</div>
	{/if}
	<h1 class="text-center font-bold max-w-lg leading-tight text-2xl md:text-3xl mt-1">
		{signInStatus === 'sent-otp'
			? $LL.SignIn.PageTitleSentLink()
			: $LL.SignIn.PageTitleCreateAccountOrSignIn()}
	</h1>
	<div class="w-full flex flex-col items-center justify-start mt-3">
		<p class="max-w-lg text-c-on-bg/60 text-center leading-relaxed mb-6">
			{signInStatus === 'sent-otp'
				? $LL.SignIn.PageParagraphSentLink()
				: $LL.SignIn.PageParagraph()}
		</p>
		{#if signInStatus === 'sent-otp'}
			<div
				transition:expandCollapse|local={{ duration: 200, easing: quadOut, opacity: 0 }}
				class="w-full flex flex-col items-center justify-start relative z-0"
			>
				<div class="w-full p-1 flex items-center justify-center">
					<NoBgButton onClick={() => (signInStatus = 'idle')}>
						<div class="flex items-center justify-center gap-2.5 px-2 py-1">
							<IconBack
								class="w-6 h-6 transform transition text-c-on-bg/50 group-hover:-translate-x-1
                {!$isTouchscreen ? 'group-hover:text-c-primary' : ''}"
							/>
							<p class="font-bold">{$LL.Shared.GoBackButton()}</p>
						</div>
					</NoBgButton>
				</div>
			</div>
		{:else}
			<div
				transition:expandCollapse|local={{ duration: 200, easing: quadOut, opacity: 0 }}
				class="w-full relative z-0 flex flex-col justify-start
        {!isModal
					? 'max-w-sm rounded-3xl bg-c-bg-secondary ring-2 ring-c-bg-tertiary shadow-xl shadow-c-shadow/[var(--o-shadow-normal)]'
					: 'max-w-xs'}"
			>
				<form
					on:input={() => (errorText = null)}
					on:submit|preventDefault={signIn}
					class="w-full flex flex-col {!isModal ? 'p-4 md:p-5' : 'p-1 md:pb-2'}"
				>
					<div class="w-full flex flex-col justify-start">
						<Input
							disabled={signInStatus === 'loading'}
							type="email"
							title={$LL.Shared.EmailInput.Placeholder()}
							bind:value={email}
							hasIcon
						>
							<IconEmail slot="icon" class="w-full h-full" />
						</Input>
					</div>
					{#if errorText}
						<ErrorLine text={errorText} class="text-xs" />
					{/if}
					<Button class="mt-4" loading={signInStatus === 'loading'} withSpinner>
						{$LL.SignIn.ContinueButton()}
					</Button>
				</form>
				<!-- <div
					class="w-full flex flex-col items-center justify-center {!isModal
						? 'p-4 md:p-5'
						: 'p-1 md:pb-2'}"
				>
					<p>Or continue with:</p>
					<button
						on:click={() => signInWithOAuth('google')}
						class="w-full px-6 md:px-8 py-5 rounded-xl bg-white text-black/75 flex items-center justify-center font-bold gap-3"
					>
						<IconGoogle class="w-6 h-6 flex-shrink-0" />
						<p>Sign in with Google</p>
					</button>
				</div> -->
			</div>
		{/if}
	</div>
</div>
