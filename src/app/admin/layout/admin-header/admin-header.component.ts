import { Component } from '@angular/core';
import { LogoComponent } from '../../../core/components/logo/logo.component';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-admin-header',
  imports: [LogoComponent, ButtonModule],
  templateUrl: './admin-header.component.html',
  styleUrl: './admin-header.component.css',
})
export class AdminHeaderComponent {
  logOut() {}
}
