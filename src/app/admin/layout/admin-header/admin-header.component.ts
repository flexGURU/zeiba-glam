import { Component } from '@angular/core';
import { LogoComponent } from '../../../core/components/logo/logo.component';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-admin-header',
  imports: [LogoComponent, ButtonModule],
  templateUrl: './admin-header.component.html',
  styleUrl: './admin-header.component.css',
})
export class AdminHeaderComponent {
  constructor(private authService: AuthService) {}
  logOut() {
    this.authService.logout();
  }
}
