/*
 *  Licensed to the Apache Software Foundation (ASF) under one
 *  or more contributor license agreements.  See the NOTICE file
 *  distributed with this work for additional information
 *  regarding copyright ownership.  The ASF licenses this file
 *  to you under the Apache License, Version 2.0 (the
 *  "License"); you may not use this file except in compliance
 *  with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

/**
 * @author Igor Ostapenko
 */
'use strict';

const assert = require('assert');
const { traverserSerializer } = require('../../../lib/structure/io/binary/GraphBinary');
const { Traverser } = require('../../../lib/process/traversal');

const { from, concat } = Buffer;

describe('GraphBinary.TraverserSerializer', () => {

  const type_code =  from([0x21]);
  const value_flag = from([0x00]);

  const cases = [
    { v:undefined,                           fq:1, b:[0x21, 0x01],                                         av:null },
    { v:undefined,                           fq:0, b:[0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x01, 0xFE,0x01], av:new Traverser(null, 1n) },
    { v:null,                                fq:1, b:[0x21, 0x01] },
    { v:null,                                fq:0, b:[0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x01, 0xFE,0x01], av:new Traverser(null, 1n) },

    { v:new Traverser(-1, 0n),                     b:[0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00, 0x01,0x00,0xFF,0xFF,0xFF,0xFF] },
    { v:new Traverser('abC', 2n),                  b:[0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x02, 0x03,0x00, 0x00,0x00,0x00,0x03, 0x61,0x62,0x43] },

    { des:1, err:/buffer is missing/,        fq:1, b:undefined },
    { des:1, err:/buffer is missing/,        fq:0, b:undefined },
    { des:1, err:/buffer is missing/,        fq:1, b:null },
    { des:1, err:/buffer is missing/,        fq:0, b:null },
    { des:1, err:/buffer is empty/,          fq:1, b:[] },
    { des:1, err:/buffer is empty/,          fq:0, b:[] },

    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0x00] },
    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0x01] },
    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0x20] },
    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0x22] },
    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0x12] },
    { des:1, err:/unexpected {type_code}/,   fq:1, b:[0xFF] },

    { des:1, err:/{value_flag} is missing/,  fq:1, b:[0x21] },
    { des:1, err:/unexpected {value_flag}/,  fq:1, b:[0x21,0x10] },
    { des:1, err:/unexpected {value_flag}/,  fq:1, b:[0x21,0x02] },
    { des:1, err:/unexpected {value_flag}/,  fq:1, b:[0x21,0x0F] },
    { des:1, err:/unexpected {value_flag}/,  fq:1, b:[0x21,0xFF] },

    { des:1, err:/{bulk} is less than zero/,       b:[0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF] }, // TODO: is it correct expectation that bulk must be >= 0?
    { des:1, err:/{bulk} is less than zero/,       b:[0x80,0x00,0x00,0x00,0x00,0x00,0x00,0x00] },
  ];

  describe.skip('serialize', () => {
    cases.forEach(({ des, v, fq, b }, i) => it(`should be able to handle case #${i}`, () => {
      // deserialize case only
      if (des)
        return; // keep it like passed test not to mess with case index

      b = from(b);

      // when fq is under control
      if (fq !== undefined) {
        assert.deepEqual( traverserSerializer.serialize(v, fq), b );
        return;
      }

      // generic case
      assert.deepEqual( traverserSerializer.serialize(v, true),  concat([type_code, value_flag, b]) );
      assert.deepEqual( traverserSerializer.serialize(v, false), concat([                       b]) );
    }));
  });

  describe('deserialize', () =>
    cases.forEach(({ v, fq, b, av, err }, i) => it(`should be able to handle case #${i}`, () => {
      if (Array.isArray(b))
        b = from(b);

      // wrong binary
      if (err !== undefined) {
        if (fq !== undefined)
          assert.throws(() => traverserSerializer.deserialize(b, fq), { message: err });
        else {
          assert.throws(() => traverserSerializer.deserialize(concat([type_code, value_flag, b]), true),  { message: err });
          assert.throws(() => traverserSerializer.deserialize(concat([                       b]), false), { message: err });
        }
        return;
      }

      if (av !== undefined)
        v = av;
      const len = b.length;

      // when fq is under control
      if (fq !== undefined) {
        assert.deepStrictEqual( traverserSerializer.deserialize(b, fq), {v,len} );
        return;
      }

      // generic case
      assert.deepStrictEqual( traverserSerializer.deserialize(concat([type_code, value_flag, b]), true),  {v,len:len+2} );
      assert.deepStrictEqual( traverserSerializer.deserialize(concat([                       b]), false), {v,len:len+0} );
    }))
  );

  describe('canBeUsedFor', () =>
    it.skip('')
  );

});