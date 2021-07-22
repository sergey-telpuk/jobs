<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Jobs\Serializer;

use Spiral\RoadRunner\Jobs\Exception\SerializationException;
use Spiral\RoadRunner\Jobs\SerializerInterface;

final class JsonSerializer implements SerializerInterface
{
    /**
     * @param array $payload
     * @return string
     * @throws SerializationException
     */
    public function serialize(array $payload): string
    {
        try {
            return \json_encode($payload, \JSON_THROW_ON_ERROR);
        } catch (\Throwable $e) {
            throw new SerializationException($e->getMessage(), (int)$e->getCode(), $e);
        }
    }

    /**
     * {@inheritDoc}
     */
    public function unserialize(string $payload): array
    {
        try {
            return (array)\json_decode($payload, true, 512, \JSON_THROW_ON_ERROR);
        } catch (\Throwable $e) {
            throw new SerializationException($e->getMessage(), (int)$e->getCode(), $e);
        }
    }
}
