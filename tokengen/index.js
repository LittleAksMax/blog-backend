import { initializeApp } from 'firebase/app';
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';
import config from './firebaseconfig.json' with { type: 'json' };

const main = async () => {
    const app = initializeApp(config.firebase);
    const auth = getAuth(app);
    const creds = config.creds;
    const user = await signInWithEmailAndPassword(auth,
        creds.email, creds.password);
    console.log(await auth.currentUser.getIdToken());
}

main();