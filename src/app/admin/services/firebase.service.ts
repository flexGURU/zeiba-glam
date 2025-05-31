import { Injectable } from '@angular/core';
import { AngularFireStorage } from '@angular/fire/compat/storage';
import { BehaviorSubject, from, Observable, of, tap } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FirebaseService {
  imageURL = new BehaviorSubject<string | null>(null);
  imageURL$ = this.imageURL.asObservable();
  constructor(private storage: AngularFireStorage) {}

  uploadImage(file: File[], formData: any) {}

  uploadImages = (selectedImages: File[]): Promise<string[]> => {
    const uploadPromises = selectedImages.map((file) => {
      const filePath = `postIMG/${Date.now()}_${file.name}`;
      const fileRef = this.storage.ref(filePath);
      const metadata = { contentType: file.type };

      return this.storage.upload(filePath, file, metadata).then(() => {
        return fileRef.getDownloadURL().toPromise(); // Convert Observable to Promise
      });
    });

    return Promise.all(uploadPromises); // Resolve when all uploads complete
  };
}
