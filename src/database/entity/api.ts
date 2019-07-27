import { Column, Entity, JoinColumn, ManyToOne } from 'typeorm';
import { ApiType } from './api-type';
import { Base } from '../base';

@Entity()
export class Api extends Base {
  @ManyToOne(type => ApiType, {
    onDelete: 'RESTRICT',
    onUpdate: 'RESTRICT',
  })
  @JoinColumn({
    name: 'type',
    referencedColumnName: 'id',
  })
  type: number;

  @Column('json', {
    comment: '接口名称',
  })
  name: any;

  @Column('varchar', {
    length: 90,
    unique: true,
    comment: '接口地址',
  })
  api: string;
}
