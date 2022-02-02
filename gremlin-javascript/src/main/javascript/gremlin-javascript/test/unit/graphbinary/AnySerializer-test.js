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
const { DataType, anySerializer } = require('../../../lib/structure/io/binary/GraphBinary');

const t = require('../../../lib/process/traversal');
const Bytecode = require('../../../lib/process/bytecode');
const { GraphTraversal } = require('../../../lib/process/graph-traversal');
const g = () => new GraphTraversal(undefined, undefined, new Bytecode());

const { from, concat } = Buffer;

describe('GraphBinary.AnySerializer', () => {

  describe('serialize', () =>
    [
      /*
        This is the most complex part of GraphBinary,
        we are expected to guess right serializer, i.e. right type.
        GraphSON implementation splits decision making among client and server sides due to JSON nature.
        But GraphBinary have to decide on a type on client side only.

        Let's follow existing structure/io/graph-serializer.GraphSON3Writer.adaptObject() logic
        of serializer/type selection and add the rest of the types.

        These test suite is the documentation and set of examples of what type is inferred
        for a given JavaScript value.
      */

      // NumberSerializer
      // { v:NaN, b:[] },
      // { v:Infinity, b:[] },
      // { v:-Infinity, b:[] },
      // TODO

      // DateSerializer,
      // TODO

      // BytecodeSerializer
      // Traversal type
      { v: g().V(),
        b: [
          DataType.BYTECODE,0x00,
          0x00,0x00,0x00,0x01, // {steps_length}
            0x00,0x00,0x00,0x01, 0x56, // V
            0x00,0x00,0x00,0x00, // (void)
          0x00,0x00,0x00,0x00, // {sources_length}
        ]
      },
      // Bytecode type
      { v: g().V().getBytecode(),
        b: [
          DataType.BYTECODE,0x00,
          0x00,0x00,0x00,0x01, // {steps_length}
            0x00,0x00,0x00,0x01, 0x56, // V
            0x00,0x00,0x00,0x00, // (void)
          0x00,0x00,0x00,0x00, // {sources_length}
        ]
      },

      // TraverserSerializer
      { v: new t.Traverser("A1", 16),
        b: [
          DataType.TRAVERSER,0x00,
          0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x10,
          DataType.STRING,0x00, 0x00,0x00,0x00,0x02, 0x41,0x31,
        ]
      },

      // TraversalStrategySerializer
      // TODO

      // PSerializer
      // TODO

      // TextPSerializer
      // TODO

      // LambdaSerializer
      // TODO

      // EnumSerializer (actually represents different enum like types)
      { v: new t.EnumValue('Barrier', 'normSack'),
        b: [ DataType.BARRIER,0x00, 0x00,0x00,0x00,0x08, ...from('normSack') ]
      },
      { v: new t.EnumValue('Cardinality', 'single'),
        b: [ DataType.CARDINALITY,0x00, 0x00,0x00,0x00,0x06, ...from('single') ]
      },
      { v: new t.EnumValue('Column', 'keys'),
        b: [ DataType.COLUMN,0x00, 0x00,0x00,0x00,0x04, ...from('keys') ]
      },
      { v: new t.EnumValue('Direction', 'OUT'),
        b: [ DataType.DIRECTION,0x00, 0x00,0x00,0x00,0x03, ...from('OUT') ]
      },
      { v: new t.EnumValue('Operator', 'addAll'),
        b: [ DataType.OPERATOR,0x00, 0x00,0x00,0x00,0x06, ...from('addAll') ]
      },
      { v: new t.EnumValue('Order', 'desc'),
        b: [ DataType.ORDER,0x00, 0x00,0x00,0x00,0x04, ...from('desc') ]
      },
      { v: new t.EnumValue('Pick', 'any'),
        b: [ DataType.PICK,0x00, 0x00,0x00,0x00,0x03, ...from('any') ]
      },
      { v: new t.EnumValue('Pop', 'first'),
        b: [ DataType.POP,0x00, 0x00,0x00,0x00,0x05, ...from('first') ]
      },
      { v: new t.EnumValue('Scope', 'local'),
        b: [ DataType.SCOPE,0x00, 0x00,0x00,0x00,0x05, ...from('local') ]
      },
      { v: new t.EnumValue('T', 'id'),
        b: [ DataType.T,0x00, 0x00,0x00,0x00,0x02, ...from('id') ]
      },

      // VertexSerializer
      // TODO

      // EdgeSerializer
      // TODO

      // LongSerializer
      // TODO +

      // ListSerializer
      // TODO +
      { v:[],  b:[DataType.LIST,0x00, 0x00,0x00,0x00,0x00] },
      { v:[1], b:[DataType.LIST,0x00, 0x00,0x00,0x00,0x01, 0x01,0x00, 0x00,0x00,0x00,0x01] },

      // SetSerializer
      // TODO

      // MapSerializer
      // TODO +

      // Default (strings / objects / ...)
      // TODO: string
      // TODO: uuid
      // TODO: unspecified null
      // TODO: leftovers
    ].forEach(({ v, b }, i) => it(`should be able to handle case #${i}`, () => {
      b = from(b);
      assert.deepEqual(anySerializer.serialize(v, true),  b);
      assert.deepEqual(anySerializer.serialize(v, false), b.slice(2));
    }))
  );

  describe('deserialize', () =>
    [
      { err:/buffer is missing/,                  b:undefined },
      { err:/buffer is missing/,                  b:undefined },
      { err:/buffer is missing/,                  b:null },
      { err:/buffer is missing/,                  b:null },
      { err:/buffer is empty/,                    b:[] },
      { err:/buffer is empty/,                    b:[] },

      { err:/unknown {type_code}/,                b:[0x2E] },
      { err:/unknown {type_code}/,                b:[0x30] },
      { err:/unknown {type_code}/,                b:[0x8F] },
      { err:/unknown {type_code}/,                b:[0xFF] },

      // this is like a register of supported type deserializers:

      // INT
      { v:null,                                   b:[0x01,0x01] },
      { v:1,                                      b:[0x01,0x00, 0x00,0x00,0x00,0x01] },

      // LONG
      { v:null,                                   b:[0x02,0x01] },
      { v:1n,                                     b:[0x02,0x00, 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x01] },

      // STRING
      { v:null,                                   b:[0x03,0x01] },
      { v:'Ab0',                                  b:[0x03,0x00, 0x00,0x00,0x00,0x03, 0x41,0x62,0x30] },

      // LIST
      { v:null,                                   b:[0x09,0x01] },
      { v:[],                                     b:[0x09,0x00, 0x00,0x00,0x00,0x00] },

      // MAP
      { v:null,                                   b:[0x0A,0x01] },
      { v:{},                                     b:[0x0A,0x00, 0x00,0x00,0x00,0x00] },

      // UUID
      { v:null,                                   b:[0x0C,0x01] },
      { v:'00010203-0405-0607-0809-0a0b0c0d0e0f', b:[0x0C,0x00, 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,0x09,0x0A,0x0B,0x0C,0x0D,0x0E,0x0F] },

      // BARRIER
      { v:new t.EnumValue('Barrier',null),        b:[0x13,0x01] },
      { v:new t.EnumValue('Barrier','normSack'),  b:[0x13,0x00, 0x00,0x00,0x00,0x08, ...from('normSack')] },

      // BYTECODE
      { v:null,                                   b:[0x15,0x01] },
      { v:new Bytecode(),                         b:[0x15,0x00, 0x00,0x00,0x00,0x00, 0x00,0x00,0x00,0x00] },

      // CARDINALITY
      { v:new t.EnumValue('Cardinality',null),    b:[0x16,0x01] },
      { v:new t.EnumValue('Cardinality', 'set'),  b:[0x16,0x00, 0x00,0x00,0x00,0x03, ...from('set')] },

      // COLUMN
      { v:new t.EnumValue('Column',null),         b:[0x17,0x01] },
      { v:new t.EnumValue('Column','keys'),       b:[0x17,0x00, 0x00,0x00,0x00,0x04, ...from('keys')] },

      // DIRECTION
      { v:new t.EnumValue('Direction',null),      b:[0x18,0x01] },
      { v:new t.EnumValue('Direction','OUT'),     b:[0x18,0x00, 0x00,0x00,0x00,0x03, ...from('OUT')] },

      // OPERATOR
      { v:new t.EnumValue('Operator',null),       b:[0x19,0x01] },
      { v:new t.EnumValue('Operator','addAll'),   b:[0x19,0x00, 0x00,0x00,0x00,0x06, ...from('addAll')] },

      // ORDER
      { v:new t.EnumValue('Order',null),          b:[0x1A,0x01] },
      { v:new t.EnumValue('Order','desc'),        b:[0x1A,0x00, 0x00,0x00,0x00,0x04, ...from('desc')] },

      // PICK
      { v:new t.EnumValue('Pick',null),           b:[0x1B,0x01] },
      { v:new t.EnumValue('Pick','any'),          b:[0x1B,0x00, 0x00,0x00,0x00,0x03, ...from('any')] },

      // POP
      { v:new t.EnumValue('Pop',null),            b:[0x1C,0x01] },
      { v:new t.EnumValue('Pop','first'),         b:[0x1C,0x00, 0x00,0x00,0x00,0x05, ...from('first')] },

      // SCOPE
      { v:new t.EnumValue('Scope',null),          b:[0x1F,0x01] },
      { v:new t.EnumValue('Scope','local'),       b:[0x1F,0x00, 0x00,0x00,0x00,0x05, ...from('local')] },

      // T
      { v:new t.EnumValue('T',null),              b:[0x20,0x01] },
      { v:new t.EnumValue('T','id'),              b:[0x20,0x00, 0x00,0x00,0x00,0x02, ...from('id')] },

      // TRAVERSER
      { v:null,                                   b:[0x21,0x01] },
      { v:new t.Traverser('A', 2n),               b:[0x21,0x00, 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x02, 0x03,0x00,0x00,0x00,0x00,0x01,0x41] },

      // UNSPECIFIED_NULL
      { v:null,                                   b:[0xFE,0x01] },

      // TODO: "register" other types
    ]
    .forEach(({ v, b, err }, i) => it(`should be able to handle case #${i}`, () => {
      if (Array.isArray(b))
        b = from(b);

      // wrong binary
      if (err !== undefined) {
        assert.throws(() => anySerializer.deserialize(b), { message: err });
        return;
      }

      const len = b.length;
      assert.deepStrictEqual( anySerializer.deserialize(b), {v,len} );
    }))
  );

});
