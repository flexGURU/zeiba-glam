import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { CheckboxModule } from 'primeng/checkbox';
import { SelectModule } from 'primeng/select';
import { ToastModule } from 'primeng/toast';
import { PasswordModule } from 'primeng/password';
import { MessageModule } from 'primeng/message';
import { AuthService } from '../../services/auth.service';
import { Router } from 'express';
import { MessageService } from 'primeng/api';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';
import { LogoComponent } from '../../../core/components/logo/logo.component';
import { User } from '../../services/user.interface';

@Component({
  selector: 'app-login',
  imports: [
    SelectModule,
    FormsModule,
    CommonModule,
    PasswordModule,
    CheckboxModule,
    ButtonModule,
    MessageModule,
    ReactiveFormsModule,
    ToastModule,
    InputTextModule,
    CardModule,
    LogoComponent,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
  providers: [MessageService],
})
export class LoginComponent {
  loginForm: FormGroup;
  authService = inject(AuthService);
  isLoading: boolean = false;

  constructor(
    private fb: FormBuilder,
    private messageService: MessageService
  ) {
    this.loginForm = this.fb.group({
      email: ['test@gmail.com', [Validators.required, Validators.email]],
      password: ['secret', Validators.required],
    });
  }

  onLogin() {
    this.isLoading = true;
    if (this.loginForm.invalid) {
      return;
    }

    const { email, password } = this.loginForm.getRawValue();
    console.log(email, password);

    this.authService.login(email, password).subscribe({
      next: (response) => {
        if (response) {
          this.isLoading = false;
          this.messageService.add({
            severity: 'success',
            summary: 'Login',
            detail: 'Success',
          });
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.messageService.add({ severity: 'error', summary: 'Error', detail: error });
      },
    });
  }
}
