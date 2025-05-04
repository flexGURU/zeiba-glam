import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'home',
    loadChildren: () =>
      import('./features/home/routes/home.routes').then((m) => m.routes),
  },
  { path: '', redirectTo: '/home', pathMatch: 'full' },
];
