import { CommonModule } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { MessageService } from 'primeng/api';
import { InputText } from 'primeng/inputtext';
import { ToastModule } from 'primeng/toast';
import { RadioButtonModule } from 'primeng/radiobutton';
import { ButtonModule } from 'primeng/button';
import { TextareaModule } from 'primeng/textarea';
import { InputNumberModule } from 'primeng/inputnumber';
import { OrderSummaryComponent } from '../order-summary/order-summary.component';
import { CartItem, CartService } from '../../../services/cart.service';

@Component({
  selector: 'app-checkout',
  imports: [
    CommonModule,
    TextareaModule,
    ButtonModule,
    FormsModule,
    ToastModule,
    InputText,
    RadioButtonModule,
    ReactiveFormsModule,
    FormsModule,
    InputNumberModule,
    OrderSummaryComponent,
  ],
  templateUrl: './checkout.component.html',
  styleUrl: './checkout.component.css',
  providers: [MessageService],
})
export class CheckoutComponent {
  processPaymentStatus: boolean = true;
  orderItem: any = null;
  paymentMethod: string = 'mpesa';
  paymentForm: FormGroup;
  cartItems: CartItem[] = [];
  totalAmount: number = 0;

  // Customer details
  customerDetails = {
    name: '',
    email: '',
    phone: '',
    address: '',
  };

  // M-Pesa details
  mpesaPhone = signal<string>('');

  // Card details
  cardDetails = {
    number: '',
    expiry: '',
    cvv: '',
    name: '',
  };

  loading: boolean = false;

  private cartService = inject(CartService);

  constructor(
    private router: Router,
    private messageService: MessageService,
    private fb: FormBuilder
  ) {
    this.paymentForm = this.fb.nonNullable.group({
      fullName: ['', [Validators.required, Validators.minLength(4)]],
      email: ['', [Validators.required, Validators.email]],
      phone: ['', [Validators.required]],
      deliveryAddress: [''],
    });
  }

  get formControls() {
    return this.paymentForm.controls;
  }

  ngOnInit() {
    this.cartItems = this.cartService.getCartItems()
    this.totalAmount = this.cartItems.reduce((acc, item) => acc + item.total, 0);

  }

  onPaymentMethodChange(method: string) {
    this.paymentMethod = method;
  }
  updateMpesaPhone(event: any) {
    const phoneValue = event.target.value;
    this.paymentForm.get('phone')?.setValue(phoneValue);
    this.mpesaPhone = phoneValue;
  }

  processPayment() {
    // Validate customer details
    if (!this.paymentForm) {
      this.messageService.add({
        severity: 'error',
        summary: 'Missing Information',
        detail: 'Please fill in all required customer details',
      });
      return;
    }
    console.log(this.paymentForm.getRawValue());

    // Validate payment method specific details
    if (this.paymentMethod === 'mpesa' && !this.mpesaPhone) {
      this.messageService.add({
        severity: 'error',
        summary: 'Missing Phone Number',
        detail: 'Please enter your M-Pesa phone number',
      });
      return;
    }

    this.loading = true;

    // Simulate payment processing
    setTimeout(() => {
      this.loading = false;

      // Show success message
      this.messageService.add({
        severity: 'success',
        summary: 'Payment Successful',
        detail: 'Your order has been placed successfully!',
      });

      // Navigate to success page or home after delay
      setTimeout(() => {
        this.paymentForm.reset();
        this.router.navigate(['/order-success']);
      }, 2000);
    }, 3000);
  }

  goBack() {
    window.history.back();
  }
}
