import { Component, Input, OnInit } from '@angular/core';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-flipping-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './flipping-view.component.html',
  styleUrl: './flipping-view.component.scss',
})
export class FlippingViewComponent implements OnInit {
  @Input() listingData$!: Observable<IItemListings[]>;

  ngOnInit(): void {}
}
