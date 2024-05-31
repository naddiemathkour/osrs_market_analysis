import { HttpClient, HttpHeaders } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../environments/environments';
import { ILatestData } from '../interfaces/latest-data-interface';

@Injectable({
  providedIn: 'root',
})
export class MarketDataService {
  headers = new HttpHeaders({
    'User-Agent': 'Runescape Market Data Analysis'
  });

  private http = inject(HttpClient);

  fetchMarketData(): Observable<ILatestData[]> {
    const url = `${environment.apiUrl}latest`;
    return this.http.get(url, {
      headers: this.headers,
    }) as any;
  }
}
