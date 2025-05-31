import { ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideClientHydration, withEventReplay } from '@angular/platform-browser';

import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';
import Aura from '@primeng/themes/aura';
import { provideHttpClient, withFetch } from '@angular/common/http';
import MyPurplePreset from './preset';

import { AngularFirestoreModule } from '@angular/fire/compat/firestore';
import { AngularFireModule } from '@angular/fire/compat';
import { environment } from '../environments/environment.development';

export const appConfig: ApplicationConfig = {
  providers: [
    importProvidersFrom([
      AngularFireModule.initializeApp(environment.firebaseConfig),
      AngularFirestoreModule,
    ]),
    provideAnimationsAsync(),
    provideRouter(routes),
    provideHttpClient(withFetch()),
    // provideClientHydration(withEventReplay()),
    providePrimeNG({
      theme: {
        preset: MyPurplePreset,
        options: {
          darkModeSelector: '.my-app-dark',
        },
      },
    }),
  ],
};
